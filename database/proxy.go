package database

import (
	"github.com/NubeDev/flow-framework/model"
)

func (d *GormDatabase) GetFN(uuid string) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	if err := d.DB.Where("uuid = ? ", uuid).First(&flowNetworkModel).Error; err != nil {
		return nil, err
	}
	return flowNetworkModel, nil
}
