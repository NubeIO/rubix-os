package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)




// GetFlowNetworks returns all networks.
func (d *GormDatabase) GetFlowNetworks(withChildren bool) ([]*model.FlowNetwork, error) {
	var flowNetworksModel []*model.FlowNetwork
	if withChildren { // drop child to reduce json size
		query := d.DB.Find(&flowNetworksModel);if query.Error != nil {
			return nil, query.Error
		}
		return flowNetworksModel, nil
	} else {
		query := d.DB.Find(&flowNetworksModel);if query.Error != nil {
			return nil, query.Error
		}
		return flowNetworksModel, nil
	}
}

// GetFlowNetwork returns the network for the given id or nil.
func (d *GormDatabase) GetFlowNetwork(uuid string) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
		query := d.DB.Where("uuid = ? ", uuid).First(&flowNetworkModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return flowNetworkModel, nil
}

// CreateFlowNetwork creates a device.
func (d *GormDatabase) CreateFlowNetwork(body *model.FlowNetwork) (*model.FlowNetwork, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.FlowNetwork)
	body.Name = nameIsNil(body.Name)
	body.GlobalFlowID = nameIsNil(body.Name)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}



// UpdateFlowNetwork returns the network for the given id or nil.
func (d *GormDatabase) UpdateFlowNetwork(uuid string, body *model.FlowNetwork) (*model.FlowNetwork, error) {
	var flowNetworkModel *model.FlowNetwork
	query := d.DB.Where("uuid = ?", uuid).Find(&flowNetworkModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&flowNetworkModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return flowNetworkModel, nil

}

// DeleteFlowNetwork delete a network.
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

// DropFlowNetworks delete all networks.
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
