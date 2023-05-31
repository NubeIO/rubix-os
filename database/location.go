package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetLocations() ([]*model.Location, error) {
	var locationsModel []*model.Location
	query := d.buildLocationQuery()
	if err := query.Find(&locationsModel).Error; err != nil {
		return nil, err
	}
	return locationsModel, nil
}

func (d *GormDatabase) GetLocation(uuid string) (*model.Location, error) {
	var locationModel *model.Location
	query := d.buildLocationQuery()
	if err := query.Where("uuid = ?", uuid).First(&locationModel).Error; err != nil {
		return nil, err
	}
	return locationModel, nil
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

func (d *GormDatabase) DeleteLocation(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.Location{})
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DropLocations() (bool, error) {
	query := d.DB.Where("1 = 1").Delete(&model.Location{})
	return d.deleteResponseBuilder(query)
}
