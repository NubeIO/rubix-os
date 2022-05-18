package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) GetStreamClones(args api.Args) ([]*model.StreamClone, error) {
	var streamClonesModel []*model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query.Find(&streamClonesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamClonesModel, nil
}

func (d *GormDatabase) GetStreamCloneByArg(args api.Args) (*model.StreamClone, error) {
	var streamClonesModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query.Find(&streamClonesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamClonesModel, nil
}

func (d *GormDatabase) GetStreamClone(uuid string, args api.Args) (*model.StreamClone, error) {
	var streamCloneModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query = query.Where("uuid = ? ", uuid).First(&streamCloneModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamCloneModel, nil
}

func (d *GormDatabase) DeleteStreamClone(uuid string) (bool, error) {
	var streamCloneModel *model.StreamClone
	query := d.DB.Where("uuid = ? ", uuid).Delete(&streamCloneModel)
	return d.deleteResponseBuilder(query)
}

func (d *GormDatabase) DeleteOneStreamCloneByArgs(args api.Args) (bool, error) {
	var streamCloneModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	if err := query.First(&streamCloneModel).Error; err != nil {
		return false, err
	}
	query = query.Delete(&streamCloneModel)
	return d.deleteResponseBuilder(query)
}
