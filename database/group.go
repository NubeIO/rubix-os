package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/src/cli/cligetter"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"sync"
)

func (d *GormDatabase) GetGroups(args api.Args) ([]*model.Group, error) {
	var groupsModel []*model.Group
	query := d.buildGroupQuery(args)
	if err := query.Find(&groupsModel).Error; err != nil {
		return nil, err
	}
	return groupsModel, nil
}

func (d *GormDatabase) GetGroup(uuid string, args api.Args) (*model.Group, error) {
	var groupModel *model.Group
	query := d.buildGroupQuery(args)
	if err := query.Where("uuid = ?", uuid).First(&groupModel).Error; err != nil {
		return nil, err
	}
	attachOpenVPN(groupModel.Hosts)
	return groupModel, nil
}

func (d *GormDatabase) GetGroupsByUUIDs(uuids []*string, args api.Args) ([]*model.Group, error) {
	var groupsModel []*model.Group
	query := d.buildGroupQuery(args)
	if err := query.Where("uuid IN ?", uuids).Find(&groupsModel).Error; err != nil {
		return nil, err
	}
	return groupsModel, nil
}

func (d *GormDatabase) GetGroupsByHostUUIDs(hostUUIDs []*string, args api.Args) ([]*model.Group, error) {
	var groupsModel []*model.Group
	query := d.buildGroupQuery(args)
	if err := query.Distinct("groups.*").
		Joins("JOIN hosts ON groups.uuid = hosts.group_uuid").
		Where("hosts.uuid IN ?", hostUUIDs).
		Find(&groupsModel).Error; err != nil {
		return nil, err
	}
	return groupsModel, nil
}

func (d *GormDatabase) CreateGroup(body *model.Group) (*model.Group, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Group)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateGroup(uuid string, body *model.Group) (*model.Group, error) {
	var groupModel *model.Group
	if err := d.DB.Where("uuid = ?", uuid).Find(&groupModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return groupModel, nil
}

func (d *GormDatabase) DeleteGroup(uuid string) (*interfaces.Message, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.Group{})
	return d.deleteResponse(query)
}

func (d *GormDatabase) DropGroups() (*interfaces.Message, error) {
	query := d.DB.Where("1 = 1").Delete(&model.Location{})
	return d.deleteResponse(query)
}

func (d *GormDatabase) UpdateHostsStatus(uuid string) (*model.Group, error) {
	groupModel := model.Group{}
	query := d.buildGroupQuery(api.Args{WithHosts: true})
	err := query.Where("uuid = ?", uuid).Find(&groupModel).Error
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	for _, host := range groupModel.Hosts {
		wg.Add(1)
		cli := cligetter.GetEdgeClientFastTimeout(host)
		go func(h *model.Host) {
			defer wg.Done()
			globalUUID, pingable, isValidToken := cli.Ping()
			if globalUUID != nil {
				h.GlobalUUID = *globalUUID
			}
			h.IsOnline = &pingable
			h.IsValidToken = &isValidToken
		}(host)
	}
	wg.Wait()
	tx := d.DB.Begin()
	for _, host := range groupModel.Hosts {
		if err := tx.Where("uuid = ?", host.UUID).Updates(&host).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &groupModel, nil
}
