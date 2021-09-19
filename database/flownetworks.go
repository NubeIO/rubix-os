package database

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)

func (d *GormDatabase) GetFlowNetworks(args api.Args) ([]*model.FlowNetwork, error) {
	var flowNetworksModel []*model.FlowNetwork
	query := d.buildFlowNetworkQuery(args)
	if err := query.Find(&flowNetworksModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworksModel, nil
}

func (d *GormDatabase) GetFlowNetwork(uuid string, args api.Args) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.buildFlowNetworkQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&flowNetworkModel).Error; err != nil {
		return nil, query.Error

	}
	return flowNetworkModel, nil
}

func (d *GormDatabase) GetOneFlowNetworkByArgs(args api.Args) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.buildFlowNetworkQuery(args)
	if err := query.First(&flowNetworkModel).Error; err != nil {
		return nil, query.Error
	}
	return flowNetworkModel, nil
}

func (d *GormDatabase) CreateFlowNetwork(body *model.FlowNetwork) (*model.FlowNetwork, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.FlowNetwork)
	body.Name = nameIsNil(body.Name)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateFlowNetwork(uuid string, body *model.FlowNetwork) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	if err := d.DB.Where("uuid = ?", uuid).Find(&flowNetworkModel).Error; err != nil {
		return nil, err
	}
	if len(body.Streams) > 0 {
		if err := d.DB.Model(&flowNetworkModel).Association("Streams").Replace(body.Streams); err != nil {
			return nil, err
		}
	}
	if err := d.DB.Model(&flowNetworkModel).Updates(body).Error; err != nil {
		return nil, err
	}
	return flowNetworkModel, nil
}

func (d *GormDatabase) DeleteFlowNetwork(uuid string) (bool, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.DB.Where("uuid = ? ", uuid).Delete(&flowNetworkModel)
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

func (d *GormDatabase) DropFlowNetworks() (bool, error) {
	var networkModel *model.FlowNetwork
	query := d.DB.Where("1 = 1").Delete(&networkModel)
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
