package database

import (
	"github.com/NubeDev/flow-framework/model"
)

var networksModel []*model.Network
var networkModel *model.Network
var deviceChildTable = "Device"

// GetNetworks returns all networks.
func (d *GormDatabase) GetNetworks() ([]*model.Network, error) {
	withChildren := true
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
func (d *GormDatabase) GetNetwork(uuid string) (*model.Network, error) {
	withChildren := false
	if withChildren { // drop child to reduce json size
		query := d.DB.Where("uuid = ? ", uuid).Preload(deviceChildTable).First(&networkModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networkModel, nil
	} else {
		query := d.DB.Where("uuid = ? ", uuid).Find(&networkModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return networkModel, nil
	}
}

// CreateNetwork creates a network.
func (d *GormDatabase) CreateNetwork(network *model.Network) error {
	n := d.DB.Create(network).Error
	return n
}





// UpdateNetwork returns the network for the given id or nil.
func (d *GormDatabase) UpdateNetwork(uuid string, body *model.Network) (*model.Network, error) {
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
