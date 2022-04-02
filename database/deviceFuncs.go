package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

//GetDeviceByPointUUID get a device by its pointUUID
func (d *GormDatabase) GetDeviceByPointUUID(point *model.Point) (*model.Device, error) {
	device, err := d.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (d *GormDatabase) GetOneDeviceByArgs(args api.Args) (*model.Device, error) {
	var deviceModel *model.Device
	query := d.buildDeviceQuery(args)
	if err := query.First(&deviceModel).Error; err != nil {
		return nil, err
	}
	return deviceModel, nil
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
