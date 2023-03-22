package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
)

func (d *GormDatabase) SyncDevice(body *interfaces.SyncDevice, network *model.Network) (*model.Device, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	device, err := d.GetDeviceByName(network.Name, body.DeviceName, api.Args{WithTags: true})
	if err != nil {
		deviceModel := &model.Device{}
		deviceModel.Name = body.DeviceName
		deviceModel.Enable = boolean.NewTrue()
		deviceModel.NetworkUUID = network.UUID
		deviceModel.CreatedFromAutoMapping = boolean.NewTrue()
		deviceModel.Tags = body.DeviceTags
		deviceModel.MetaTags = body.DeviceMetaTags
		device, err = d.CreateDevice(deviceModel)
		return device, err
	}
	_, _ = d.CreateDeviceMetaTags(device.UUID, body.DeviceMetaTags)
	if !reflect.DeepEqual(device.Tags, body.DeviceTags) {
		_ = d.updateTags(&device, body.DeviceTags)
	}
	return device, nil
}
