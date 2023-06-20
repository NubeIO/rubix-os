package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetViewSetting() (*model.ViewSetting, error) {
	var viewSettingModel *model.ViewSetting
	if err := d.DB.First(&viewSettingModel).Error; err != nil {
		return nil, err
	}
	return viewSettingModel, nil
}

func (d *GormDatabase) UpsertSetting(body *model.ViewSetting) (*model.ViewSetting, error) {
	viewSetting, _ := d.GetViewSetting()
	if viewSetting != nil {
		if err := d.DB.Model(&viewSetting).Updates(body).Error; err != nil {
			return nil, err
		}
		return viewSetting, nil
	}
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.ViewSetting)
	body.Logo = marshalJson(body.Logo)
	body.Theme = marshalJson(body.Theme)
	body.WidgetConfig = marshalJson(body.WidgetConfig)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) DeleteViewSetting() (bool, error) {
	query := d.DB.Where("1 = 1").Delete(&model.ViewSetting{})
	return d.deleteResponseBuilder(query)
}
