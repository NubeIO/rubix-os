package database

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
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

func (d *GormDatabase) GetStreamClone(uuid string, args api.Args) (*model.StreamClone, error) {
	var streamCloneModel *model.StreamClone
	query := d.buildStreamCloneQuery(args)
	query = query.Where("uuid = ? ", uuid).First(&streamCloneModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return streamCloneModel, nil
}
