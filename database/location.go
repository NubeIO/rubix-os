package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	parentArgs "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetLocations(args parentArgs.Args) ([]*model.Location, error) {
	var locationsModel []*model.Location
	query := d.buildLocationQuery(args)
	if err := query.Find(&locationsModel).Error; err != nil {
		return nil, err
	}
	return locationsModel, nil
}

func (d *GormDatabase) GetLocation(uuid string, args parentArgs.Args) (*model.Location, error) {
	var locationModel *model.Location
	query := d.buildLocationQuery(args)
	if err := query.Where("uuid = ?", uuid).First(&locationModel).Error; err != nil {
		return nil, err
	}
	return locationModel, nil
}

func (d *GormDatabase) GetLocationsByUUIDs(uuids []*string, args parentArgs.Args) ([]*model.Location, error) {
	var locationsModel []*model.Location
	query := d.buildLocationQuery(args)
	if err := query.Where("uuid IN ?", uuids).Find(&locationsModel).Error; err != nil {
		return nil, err
	}
	return locationsModel, nil
}

func (d *GormDatabase) GetLocationsByGroupAndHostUUIDs(groupUUIDs []*string, hostUUIDs []*string) ([]*model.Location,
	error) {
	var locationsModel []*model.Location
	query := d.buildGroupQuery(parentArgs.Args{WithViews: true})
	query = query.Distinct("locations.*").
		Joins("JOIN groups ON locations.uuid = groups.location_uuid").
		Joins("JOIN hosts ON groups.uuid = hosts.group_uuid").
		Where("groups.uuid IN ?", groupUUIDs).Or("hosts.uuid IN ?", hostUUIDs)
	if err := query.Find(&locationsModel).Error; err != nil {
		return nil, err
	}
	return locationsModel, nil
}

func (d *GormDatabase) GetLocationGroupHostNamesByHostUUID(hostUUID string) (*interfaces.LocationGroupHostName, error) {
	var locationGroupHostNameModel *interfaces.LocationGroupHostName
	query := d.DB.Distinct("locations.name AS location_name,groups.name as group_name,hosts.name AS host_name").
		Table("locations").
		Joins("JOIN groups ON locations.uuid = groups.location_uuid").
		Joins("JOIN hosts ON groups.uuid = hosts.group_uuid").
		Where("hosts.uuid = ?", hostUUID)
	if err := query.Scan(&locationGroupHostNameModel).Error; err != nil {
		return nil, err
	}
	return locationGroupHostNameModel, nil
}

func (d *GormDatabase) CreateLocation(body *model.Location) (*model.Location, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Location)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateLocation(uuid string, body *model.Location) (*model.Location, error) {
	var locationModel *model.Location
	if err := d.DB.Where("uuid = ?", uuid).Find(&locationModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return locationModel, nil
}

func (d *GormDatabase) DeleteLocation(uuid string) (*interfaces.Message, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.Location{})
	return d.deleteResponse(query)
}

func (d *GormDatabase) DropLocations() (*interfaces.Message, error) {
	query := d.DB.Where("1 = 1").Delete(&model.Location{})
	return d.deleteResponse(query)
}
