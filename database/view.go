package database

import (
	"encoding/json"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetView(uuid string) (*model.View, error) {
	var viewModel *model.View
	if err := d.DB.Where("uuid = ?", uuid).First(&viewModel).Error; err != nil {
		return nil, err
	}
	return viewModel, nil
}

func (d *GormDatabase) GetViewsByUUIDs(uuids []*string) ([]*model.View, error) {
	var viewsModel []*model.View
	if err := d.DB.Where("uuid IN ?", uuids).Find(&viewsModel).Error; err != nil {
		return nil, err
	}
	return viewsModel, nil
}

func (d *GormDatabase) CreateView(body *model.View) (*model.View, error) {
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.View)
	body.Name = name
	mLayout, err := json.Marshal(body.Layout)
	if err != nil {
		return nil, err
	}
	body.Layout = mLayout
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateView(uuid string, body *model.View) (*model.View, error) {
	name, err := validateName(body.Name)
	if err != nil {
		return nil, err
	}
	var viewModel *model.View
	mLayout, err := json.Marshal(body.Layout)
	if err != nil {
		return nil, err
	}
	body.Layout = mLayout
	body.Name = name
	if err := d.DB.Where("uuid = ?", uuid).Find(&viewModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return viewModel, nil
}

func (d *GormDatabase) DeleteView(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.View{})
	return d.deleteResponseBuilder(query)
}
