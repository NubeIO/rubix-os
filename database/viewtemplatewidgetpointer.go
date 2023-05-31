package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetViewTemplateWidgetPointers() ([]*model.ViewTemplateWidgetPointer, error) {
	var viewTemplateWidgetPointersModel []*model.ViewTemplateWidgetPointer
	if err := d.DB.Find(&viewTemplateWidgetPointersModel).Error; err != nil {
		return nil, err
	}
	return viewTemplateWidgetPointersModel, nil
}

func (d *GormDatabase) GetViewTemplateWidgetPointer(uuid string) (*model.ViewTemplateWidgetPointer, error) {
	var viewTemplateWidgetPointerModel *model.ViewTemplateWidgetPointer
	if err := d.DB.Where("uuid = ?", uuid).First(&viewTemplateWidgetPointerModel).Error; err != nil {
		return nil, err
	}
	return viewTemplateWidgetPointerModel, nil
}

func (d *GormDatabase) GetViewTemplateWidgetByViewUUID(viewUUID string) ([]*model.ViewTemplateWidgetPointer, error) {
	var viewTemplateWidgetPointersModel []*model.ViewTemplateWidgetPointer
	if err := d.DB.Where("view_uuid = ?", viewUUID).Find(&viewTemplateWidgetPointersModel).Error; err != nil {
		return nil, err
	}
	return viewTemplateWidgetPointersModel, nil
}

func (d *GormDatabase) CreateViewTemplateWidgetPointer(body *model.ViewTemplateWidgetPointer) (
	*model.ViewTemplateWidgetPointer, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.ViewTemplateWidget)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateViewTemplateWidgetPointer(uuid string, body *model.ViewTemplateWidgetPointer) (
	*model.ViewTemplateWidgetPointer, error) {
	var viewTemplateWidgetPointerModel *model.ViewTemplateWidgetPointer
	if err := d.DB.Where("uuid = ?", uuid).Find(&viewTemplateWidgetPointerModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return viewTemplateWidgetPointerModel, nil
}

func (d *GormDatabase) DeleteViewTemplateWidgetPointer(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.ViewTemplateWidgetPointer{})
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteViewTemplateWidgetPointerByViewUUID(viewUUID string) (bool, error) {
	query := d.DB.Where("view_uuid = ?", viewUUID).Delete(&model.ViewTemplateWidgetPointer{})
	return d.deleteResponseBuilder(query)
}
