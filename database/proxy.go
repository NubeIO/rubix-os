package database

import (
	"github.com/NubeIO/flow-framework/model"
)

func (d *GormDatabase) GetFN(uuid string) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	if err := d.DB.Where("uuid = ? ", uuid).First(&flowNetworkModel).Error; err != nil {
		return nil, err
	}
	return flowNetworkModel, nil
}

func (d *GormDatabase) GetFNC(uuid string) (*model.FlowNetworkClone, error) {
	var flowNetworkCloneModel *model.FlowNetworkClone
	if err := d.DB.Where("uuid = ? ", uuid).First(&flowNetworkCloneModel).Error; err != nil {
		return nil, err
	}
	return flowNetworkCloneModel, nil
}
