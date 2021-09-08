package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


var pointChildTable = "Point"

// GetDevices returns all devices.
func (d *GormDatabase) GetDevices(withPoints bool) ([]*model.Device, error) {
	var devicesModel []*model.Device
	withChildren := withPoints
	if withChildren { // drop child to reduce json size
		query := d.DB.Preload("Point").Find(&devicesModel);if query.Error != nil {
			return nil, query.Error
		}
		return devicesModel, nil
	} else {
		query := d.DB.Find(&devicesModel);if query.Error != nil {
			return nil, query.Error
		}
		return devicesModel, nil
	}
}

// GetDevice returns the device for the given id or nil.
func (d *GormDatabase) GetDevice(uuid string, withPoints bool) (*model.Device, error) {
	var deviceModel *model.Device
	withChildren := withPoints
	if withChildren { // drop child to reduce json size
		query := d.DB.Where("uuid = ? ", uuid).Preload(pointChildTable).First(&deviceModel);if query.Error != nil {
			return nil, query.Error
		}
		return deviceModel, nil
	} else {
		query := d.DB.Where("uuid = ? ", uuid).First(&deviceModel); if query.Error != nil {
			return nil, query.Error
		}
		return deviceModel, nil
	}
}

// GetDeviceByField returns the device for the given field ie name or nil.
func (d *GormDatabase) GetDeviceByField(field string, value string, withPoints bool) (*model.Device, error) {
	var deviceModel *model.Device
	f := fmt.Sprintf("%s = ? ", field)
	withChildren := withPoints
	if withChildren { // drop child to reduce json size
		query := d.DB.Where(f, value).Preload(pointChildTable).First(&deviceModel);if query.Error != nil {
			return nil, query.Error
		}
		return deviceModel, nil
	} else {
		query := d.DB.Where(f, value).First(&deviceModel); if query.Error != nil {
			return nil, query.Error
		}
		return deviceModel, nil
	}
}

// CreateDevice creates a device.
func (d *GormDatabase) CreateDevice(body *model.Device) (*model.Device, error) {
	var net *model.Network
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Device)
	networkUUID := body.NetworkUUID
	body.Name = nameIsNil(body.Name)
	query := d.DB.Where("uuid = ? ", networkUUID).First(&net);if query.Error != nil {
		return nil, query.Error
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, query.Error
	}
	return body, query.Error
}


// UpdateDevice returns the device for the given id or nil.
func (d *GormDatabase) UpdateDevice(uuid string, body *model.Device) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.DB.Where("uuid = ?", uuid).Find(&deviceModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&deviceModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return deviceModel, nil
}


// UpdateDeviceByField get by field and update.
func (d *GormDatabase) UpdateDeviceByField(field string, value string, body *model.Device) (*model.Device, error) {
	var deviceModel *model.Device
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).Find(&deviceModel);if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&deviceModel).Updates(body);if query.Error != nil {
		return nil, query.Error
	}
	return deviceModel, nil
}

// DeleteDevice delete a Device.
func (d *GormDatabase) DeleteDevice(uuid string) (bool, error) {
	var deviceModel *model.Device
	query := d.DB.Where("uuid = ? ", uuid).Delete(&deviceModel);if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}

}


// DropDevices delete all devices.
func (d *GormDatabase) DropDevices() (bool, error) {
	var deviceModel *model.Device
	query := d.DB.Where("1 = 1").Delete(&deviceModel)
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