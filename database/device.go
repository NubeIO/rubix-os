package database

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"time"
)

func (d *GormDatabase) GetDevices(args api.Args) ([]*model.Device, error) {
	var devicesModel []*model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Find(&devicesModel).Error; err != nil {
		return nil, err
	}
	return devicesModel, nil
}

func (d *GormDatabase) GetDevice(uuid string, args api.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&deviceModel).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
}

// GetDeviceByField returns the device for the given field ie name or nil.
func (d *GormDatabase) GetDeviceByField(field string, value string, withPoints bool) (*model.Device, error) {
	var deviceModel *model.Device
	f := fmt.Sprintf("%s = ? ", field)
	withChildren := withPoints
	if withChildren { // drop child to reduce json size
		query := d.DB.Where(f, value).Preload("Points").First(&deviceModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return deviceModel, nil
	} else {
		query := d.DB.Where(f, value).First(&deviceModel)
		if query.Error != nil {
			return nil, query.Error
		}
		return deviceModel, nil
	}
}

// GetPluginIDFromDevice returns the pluginUUID by using the deviceUUID to query the network.
func (d *GormDatabase) GetPluginIDFromDevice(uuid string) (*model.Network, error) {
	device, err := d.GetDevice(uuid, api.Args{})
	if err != nil {
		return nil, err
	}
	network, err := d.GetNetwork(device.NetworkUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	return network, err
}

func (d *GormDatabase) CreateDevice(body *model.Device) (*model.Device, error) {
	var net *model.Network
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Device)
	networkUUID := body.NetworkUUID
	body.Name = nameIsNil(body.Name)
	_, err := checkTransport(body.TransportType)
	if err != nil {
		return nil, err
	}
	query := d.DB.Where("uuid = ? ", networkUUID).First(&net)
	if query.Error != nil {
		return nil, query.Error
	}
	body.ThingClass = model.ThingClass.Device
	body.CommonEnable.Enable = utils.NewTrue()
	body.CommonFault.InFault = true
	body.CommonFault.MessageLevel = model.MessageLevel.NoneCritical
	body.CommonFault.MessageCode = model.CommonFaultCode.PluginNotEnabled
	body.CommonFault.Message = model.CommonFaultMessage.PluginNotEnabled
	body.CommonFault.LastFail = time.Now().UTC()
	body.CommonFault.LastOk = time.Now().UTC()
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, query.Error
	}
	var nModel *model.Network
	query = d.DB.Where("uuid = ?", body.NetworkUUID).Find(&nModel)
	if query.Error != nil {
		return nil, query.Error
	}
	t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsCreated, nModel.PluginConfId, body.UUID)
	d.Bus.RegisterTopic(t)
	err = d.Bus.Emit(eventbus.CTX(), t, body)
	if err != nil {
		return nil, errors.New("error on device eventbus")
	}
	return body, query.Error
}

func (d *GormDatabase) UpdateDevice(uuid string, body *model.Device) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.DB.Where("uuid = ?", uuid).Find(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&deviceModel).Updates(body)
	_, err := checkTransport(body.TransportType)
	if err != nil {
		return nil, err
	}
	if query.Error != nil {
		return nil, query.Error
	}
	var nModel *model.Network
	query = d.DB.Where("uuid = ?", deviceModel.NetworkUUID).Find(&nModel)
	if query.Error != nil {
		return nil, query.Error
	}
	t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsUpdated, nModel.PluginConfId, uuid)
	d.Bus.RegisterTopic(t)
	err = d.Bus.Emit(eventbus.CTX(), t, deviceModel)
	if err != nil {
		return nil, errors.New("error on device eventbus")
	}
	return deviceModel, nil
}

// UpdateDeviceByField get by field and update.
func (d *GormDatabase) UpdateDeviceByField(field string, value string, body *model.Device) (*model.Device, error) {
	var deviceModel *model.Device
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).Find(&deviceModel)
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
	var deviceModel *model.Device
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
