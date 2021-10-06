package database

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
)

func (d *GormDatabase) GetFlowNetworkClones(args api.Args) ([]*model.FlowNetworkClone, error) {
	var flowNetworkClonesModel []*model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Find(&flowNetworkClonesModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworkClonesModel, nil
}

func (d *GormDatabase) GetFlowNetworkClone(uuid string, args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&flowNetworkCloneModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworkCloneModel, nil
}

func (d *GormDatabase) GetOneFlowNetworkCloneByArgs(args api.Args) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	query := d.buildFlowNetworkCloneQuery(args)
	if err := query.First(&flowNetworkCloneModel).Error; err != nil {
		return nil, query.Error
	}
	return flowNetworkCloneModel, nil
}
