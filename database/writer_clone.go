package database

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm"
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

func (d *GormDatabase) GetOneWriterCloneByArgsTransaction(db *gorm.DB, args api.Args) (*model.WriterClone, error) {
	var wcm *model.WriterClone
	query := buildWriterCloneQueryTransaction(db, args)
	if err := query.First(&wcm).Error; err != nil {
		return nil, err
	}
	return wcm, nil
}

func (d *GormDatabase) GetOneWriterCloneByArgs(args api.Args) (*model.WriterClone, error) {
	return d.GetOneWriterCloneByArgsTransaction(d.DB, args)
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
	wc, err := d.GetWriterClone(uuid)
	if err != nil {
		return false, err
	}
	if boolean.IsTrue(wc.CreatedFromAutoMapping) {
		return false, errors.New("can't delete auto-mapped writer clone")
	}
	query := d.DB.Where("uuid = ? ", uuid).Delete(&wc)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) updateWriterClone(uuid string, body *model.WriterClone) error {
	query := d.DB.Where("uuid = ?", uuid).Updates(body)
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (d *GormDatabase) DeleteOneWriterCloneByArgs(args api.Args) (bool, error) {
	var wcm *model.WriterClone
	query := d.buildWriterCloneQuery(args).Delete(&wcm)
	return d.deleteResponseBuilder(query)
}
