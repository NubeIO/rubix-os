package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/helpers"
	"github.com/NubeDev/flow-framework/model"
)


var devicesModel []*model.Device
var deviceModel *model.Device
var pointChildTable = "Point"


// GetDevices returns all devices.
func (d *GormDatabase) GetDevices() ([]*model.Device, error) {
	withChildren := false
	if withChildren { // drop child to reduce json size
		query := d.DB.Preload("Point").Find(&devicesModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return devicesModel, nil
	} else {
		query := d.DB.Find(&devicesModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return devicesModel, nil
	}

}

// GetDevice returns the device for the given id or nil.
func (d *GormDatabase) GetDevice(uuid string) (*model.Device, error) {
	withChildren := false
	if withChildren { // drop child to reduce json size
		query := d.DB.Where("uuid = ? ", uuid).Preload(pointChildTable).First(&deviceModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return deviceModel, nil
	} else {
		query := d.DB.Where("uuid = ? ", uuid).Find(&deviceModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return deviceModel, nil
	}
}

// CreateDevice creates a device.
func (d *GormDatabase) CreateDevice(device *model.Device, body *model.Device) error {
	//c := d.DB.Create(device).Error
	fmt.Println(22222)
	body.Uuid, _ = helpers.MakeUUID()
	networkUUID := body.NetworkUuid
	fmt.Println(body.Uuid, networkUUID)
	query := d.DB.Where("uuid = ? ", networkUUID).First(&networkModel);if query.Error != nil {
		return query.Error
	}
	fmt.Println(333333)
	if err := d.DB.Create(&body).Error; err != nil {
		return query.Error
	}
	return query.Error
}


// UpdateDevice returns the device for the given id or nil.
func (d *GormDatabase) UpdateDevice(uuid string, body *model.Device) (*model.Device, error) {
	query := d.DB.Where("uuid = ?", uuid).Find(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&deviceModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return deviceModel, nil

}

// DeleteDevice delete a Device.
func (d *GormDatabase) DeleteDevice(uuid string) (bool, error) {
	query := d.DB.Where("uuid = ? ", uuid).Delete(&deviceModel)
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
