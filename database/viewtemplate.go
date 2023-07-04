package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetViewTemplates(args argspkg.Args) ([]*model.ViewTemplate, error) {
	var viewTemplatesModel []*model.ViewTemplate
	query := d.buildViewTemplateQuery(args)
	if err := query.Find(&viewTemplatesModel).Error; err != nil {
		return nil, err
	}
	return viewTemplatesModel, nil
}

func (d *GormDatabase) GetViewTemplate(uuid string, args argspkg.Args) (*model.ViewTemplate, error) {
	var viewTemplateModel *model.ViewTemplate
	query := d.buildViewTemplateQuery(args)
	if err := query.Where("uuid = ?", uuid).First(&viewTemplateModel).Error; err != nil {
		return nil, err
	}
	return viewTemplateModel, nil
}

func (d *GormDatabase) CreateViewTemplate(body *model.ViewTemplate) (*model.ViewTemplate, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.ViewTemplate)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateViewTemplate(uuid string, body *model.ViewTemplate) (*model.ViewTemplate, error) {
	var viewTemplateModel *model.ViewTemplate
	if err := d.DB.Where("uuid = ?", uuid).Find(&viewTemplateModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return viewTemplateModel, nil
}

func (d *GormDatabase) DeleteViewTemplate(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ?", uuid).Delete(&model.ViewTemplate{})
	return d.deleteResponseBuilder(query)
}
