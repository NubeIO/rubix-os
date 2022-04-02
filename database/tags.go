package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) GetTags(args api.Args) ([]*model.Tag, error) {
	var tagsModel []*model.Tag
	query := d.buildTagQuery(args)
	query.Find(&tagsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return tagsModel, nil
}

func (d *GormDatabase) GetTag(tag string, args api.Args) (*model.Tag, error) {
	var tagModel *model.Tag
	query := d.buildTagQuery(args)
	query = query.Where("tag = ? ", tag).First(&tagModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return tagModel, nil
}

func (d *GormDatabase) CreateTag(body *model.Tag) (*model.Tag, error) {
	err := d.DB.Create(&body).Error
	if err != nil {
		return nil, newDetailedError("CreateTag", "error on trying to add a new tag", err)
	}
	return body, nil
}

func (d *GormDatabase) DeleteTag(tag string) (bool, error) {
	var tagModel *model.Tag
	query := d.DB.Where("tag = ? ", tag).Delete(&tagModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
