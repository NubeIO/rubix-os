package database

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


var deviceChildTable = "Device"

// GetNetworks returns all networks.
func (d *GormDatabase) GetNetworks(withChildren bool, withPoints bool) ([]*model.Network, error) {
	var networksModel []*model.Network
	if withChildren { // drop child to reduce json size
		query := d.DB.Preload("Device").Find(&networksModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networksModel, nil
	} else {
		query := d.DB.Find(&networksModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networksModel, nil
	}

}

// GetNetwork returns the network for the given id or nil.
func (d *GormDatabase) GetNetwork(uuid string, withChildren bool, withPoints bool) (*model.Network, error) {
	var networkModel *model.Network
	if withChildren { // drop child to reduce json size
		query := d.DB.Where("uuid = ? ", uuid).Preload(deviceChildTable).First(&networkModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networkModel, nil
	} else {
		query := d.DB.Where("uuid = ? ", uuid).First(&networkModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networkModel, nil
	}
}


// CreateNetwork creates a device.
func (d *GormDatabase) CreateNetwork(body *model.Network) (*model.Network, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Network)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}



// UpdateNetwork returns the network for the given id or nil.
func (d *GormDatabase) UpdateNetwork(uuid string, body *model.Network) (*model.Network, error) {
	var networkModel *model.Network
	query := d.DB.Where("uuid = ?", uuid).Find(&networkModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&networkModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return networkModel, nil

}

// DeleteNetwork delete a network.
func (d *GormDatabase) DeleteNetwork(uuid string) (bool, error) {
	var networkModel *model.Network
	query := d.DB.Where("uuid = ? ", uuid).Delete(&networkModel)
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

// DropNetworks delete all networks.
func (d *GormDatabase) DropNetworks() (bool, error) {
	var networkModel *model.Network
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
