package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/security"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	parentArgs "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetMembers(args parentArgs.Args) ([]*model.Member, error) {
	var membersModel []*model.Member
	query := d.buildMemberQuery(args)
	query.Find(&membersModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return membersModel, nil
}

func (d *GormDatabase) GetMember(uuid string, args parentArgs.Args) (*model.Member, error) {
	var memberModel *model.Member
	query := d.buildMemberQuery(args)
	query = query.Where("uuid = ? ", uuid).First(&memberModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return memberModel, nil
}

func (d *GormDatabase) GetMemberByUsername(username string, args parentArgs.Args) (*model.Member, error) {
	var memberModel *model.Member
	query := d.buildMemberQuery(args)
	query = query.Where("username = ? ", username).First(&memberModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return memberModel, nil
}

func (d *GormDatabase) GetMemberByEmail(email string, args parentArgs.Args) (*model.Member, error) {
	var memberModel *model.Member
	query := d.buildMemberQuery(args)
	query = query.Where("email = ? ", email).First(&memberModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return memberModel, nil
}

func (d *GormDatabase) GetMembersByUUIDs(uuids []*string, args parentArgs.Args) ([]*model.Member, error) {
	var membersModel []*model.Member
	query := d.buildMemberQuery(args)
	if err := query.Where("uuid IN ?", uuids).Find(&membersModel).Error; err != nil {
		return nil, err
	}
	return membersModel, nil
}

func (d *GormDatabase) CreateMember(body *model.Member) (*model.Member, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Member)
	hashedPassword, err := security.GeneratePasswordHash(body.Password)
	if err != nil {
		return nil, err
	}
	body.Password = hashedPassword
	body.State = nstring.New(string(model.UnVerified))
	body.Permission = nstring.New(string(model.Read))
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateMember(uuid string, body *model.Member) (*model.Member, error) {
	if body.State != nil {
		obj, err := checkMemberState(*body.State)
		if err != nil {
			return nil, err
		}
		body.State = nstring.New(string(obj))
	}
	if body.Permission != nil {
		obj, err := checkMemberPermission(*body.Permission)
		if err != nil {
			return nil, err
		}
		body.Permission = nstring.New(string(obj))
	}
	var memberModel *model.Member
	query := d.DB.Where("uuid = ?", uuid).First(&memberModel)
	if query.Error != nil {
		return nil, query.Error
	}
	body.Password = memberModel.Password
	query = d.DB.Model(&memberModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return memberModel, nil
}

func (d *GormDatabase) DeleteMember(uuid string) (bool, error) {
	var memberModel *model.Member
	query := d.DB.Where("uuid = ? ", uuid).Delete(&memberModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteMemberByUsername(username string) (bool, error) {
	var memberModel *model.Member
	query := d.DB.Where("username = ? ", username).Delete(&memberModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) ChangeMemberPassword(uuid string, password string) (*interfaces.Message, error) {
	var memberModel *model.Member
	query := d.DB.Where("uuid = ?", uuid).First(&memberModel)
	if query.Error != nil {
		return nil, query.Error
	}
	hashedPassword, err := security.GeneratePasswordHash(password)
	if err != nil {
		return nil, err
	}
	query = d.DB.Model(&memberModel).Update("password", hashedPassword)
	if query.Error != nil {
		return nil, query.Error
	}
	return &interfaces.Message{Message: "your password has been changed successfully"}, nil
}

func (d *GormDatabase) GetMemberSidebars(username string, includeWithoutViews bool) ([]*model.Location, error) {
	views, err := d.GetViewsByMemberUsername(username)
	if err != nil {
		return nil, err
	}
	viewUUIDs, locationsUUIDs, groupUUIDs, hostUUIDs := getViewsUUIDs(views)

	locations, _ := d.GetLocationsByUUIDs(locationsUUIDs, parentArgs.Args{WithGroups: true, WithHosts: true, WithViews: true})

	// Remove groupUUIDs and hostUUIDs that are already covered by locations
	for _, location := range locations {
		for _, group := range location.Groups {
			groupUUIDs = filterOutItem(groupUUIDs, nstring.New(group.UUID))
			for _, host := range group.Hosts {
				hostUUIDs = filterOutItem(hostUUIDs, nstring.New(host.UUID))
			}
		}
	}

	groupLocations, _ := d.GetLocationsByGroupAndHostUUIDs(groupUUIDs, hostUUIDs)
	if locations != nil {
		locations = append(locations, groupLocations...)
	}
	groups, _ := d.GetGroupsByUUIDs(groupUUIDs, parentArgs.Args{WithViews: true, WithHosts: true})
	if groups != nil {
		// Remove hostUUIDs that are already covered by groups
		for _, group := range groups {
			for _, host := range group.Hosts {
				hostUUIDs = filterOutItem(hostUUIDs, nstring.New(host.UUID))
			}
		}
	}

	hostGroups, _ := d.GetGroupsByHostUUIDs(hostUUIDs, parentArgs.Args{WithViews: true})
	if hostGroups != nil {
		groups = append(groups, hostGroups...)
	}

	hosts, _ := d.GetHostsByUUIDs(hostUUIDs, parentArgs.Args{WithTags: true, WithComments: true, WithViews: true})
	if hosts != nil {
		// Update the relationships between hosts and groups, and groups and locations
		for _, host := range hosts {
			for i, group := range groups {
				if group.UUID == host.GroupUUID {
					groups[i].Hosts = append(groups[i].Hosts, host)
				}
			}
		}
	}

	for _, group := range groups {
		for i, location := range locations {
			if location.UUID == group.LocationUUID {
				locations[i].Groups = append(locations[i].Groups, group)
			}
		}
	}

	for _, location := range locations {
		location.Views = filterViewsByViewUUIDs(location.Views, viewUUIDs)
		// Filter groups
		var updatedGroups []*model.Group
		for _, group := range location.Groups {
			group.Views = filterViewsByViewUUIDs(group.Views, viewUUIDs)
			// Filter hosts
			var updatedHosts []*model.Host
			for _, host := range group.Hosts {
				host.Views = filterViewsByViewUUIDs(host.Views, viewUUIDs)
				if includeWithoutViews || len(host.Views) != 0 {
					updatedHosts = append(updatedHosts, host)
				}
			}
			// Update hosts
			group.Hosts = updatedHosts
			if includeWithoutViews || len(group.Views) != 0 || len(group.Hosts) != 0 {
				updatedGroups = append(updatedGroups, group)
			}
		}
		// Update groups
		location.Groups = updatedGroups
	}
	return locations, nil
}

func getViewsUUIDs(views []*model.View) (viewUUIDs []string, locationsUUIDs []*string, groupUUIDs []*string,
	hostUUIDs []*string) {
	for _, view := range views {
		viewUUIDs = append(viewUUIDs, view.UUID)
		if view.LocationUUID != nil {
			locationsUUIDs = append(locationsUUIDs, view.LocationUUID)
		}
		if view.GroupUUID != nil {
			groupUUIDs = append(groupUUIDs, view.GroupUUID)
		}
		if view.HostUUID != nil {
			hostUUIDs = append(hostUUIDs, view.HostUUID)
		}
	}
	return viewUUIDs, locationsUUIDs, groupUUIDs, hostUUIDs
}

func filterViewsByViewUUIDs(views []*model.View, viewUUIDs []string) []*model.View {
	var filteredViews []*model.View
	for _, view := range views {
		if contains(viewUUIDs, view.UUID) {
			filteredViews = append(filteredViews, view)
		}
	}
	return filteredViews
}
