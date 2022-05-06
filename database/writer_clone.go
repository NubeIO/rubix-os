package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type WriterClone struct {
	*model.WriterClone
}

func (d *GormDatabase) GetWriterClones(args api.Args) ([]*model.WriterClone, error) {
	var writerClones []*model.WriterClone
	query := d.buildWriterCloneQuery(args)
	err := query.Find(&writerClones).Error
	if err != nil {
		return nil, query.Error
	}
	return writerClones, nil
}

func (d *GormDatabase) GetWriterClone(uuid string) (*model.WriterClone, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("uuid = ? ", uuid).First(&wcm)
	if query.Error != nil {
		return nil, query.Error
	}
	return wcm, nil
}

func (d *GormDatabase) GetOneWriterCloneByArgs(args api.Args) (*model.WriterClone, error) {
	var wcm *model.WriterClone
	query := d.buildWriterCloneQuery(args)
	if err := query.First(&wcm).Error; err != nil {
		return nil, query.Error
	}
	return wcm, nil
}

func (d *GormDatabase) CreateWriterClone(body *model.WriterClone) (*model.WriterClone, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.WriterClone)
	query := d.DB.Create(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return body, nil
}

func (d *GormDatabase) DeleteWriterClone(uuid string) (bool, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("uuid = ? ", uuid).Delete(&wcm)
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

func (d *GormDatabase) DropWriterClone() (bool, error) {
	var wcm *model.WriterClone
	query := d.DB.Where("1 = 1").Delete(&wcm)
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

func (d *GormDatabase) UpdateWriterClone(writerClone *model.WriterClone, body *model.WriterClone) error {
	query := d.DB.Model(&writerClone).Updates(body)
	if query.Error != nil {
		return query.Error
	}
	return nil
}
