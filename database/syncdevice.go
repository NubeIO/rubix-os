package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
)

func (d *GormDatabase) SyncDevice(body *interfaces.SyncDevice) (*model.Network, *model.Device, error) {
	syncNetwork := &interfaces.SyncNetwork{NetworkGlobalUUID: body.NetworkGlobalUUID, NetworkName: body.NetworkName,
		NetworkTags: body.NetworkTags, NetworkMetaTags: body.NetworkMetaTags, FlowNetworkUUID: body.FlowNetworkUUID,
		IsLocal: body.IsLocal}
	network, err := d.SyncNetwork(syncNetwork)
	if err != nil {
		return nil, nil, err
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	device, err := d.GetDeviceByName(network.Name, body.DeviceName, api.Args{WithTags: true})
	if err != nil {
		fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{SourceUUID: nils.NewString(body.FlowNetworkUUID)})
		if err != nil {
			return nil, nil, err
		}
		deviceModel := &model.Device{}
		deviceModel.Name = body.DeviceName
		deviceModel.Enable = boolean.NewTrue()
		deviceModel.NetworkUUID = network.UUID
		deviceModel.AutoMappingFlowNetworkName = fnc.Name
		deviceModel.CreatedFromAutoMapping = boolean.NewTrue()
		deviceModel.Tags = body.DeviceTags
		deviceModel.MetaTags = body.DeviceMetaTags
		device, err = d.CreateDevice(deviceModel)
		return network, device, err
	}
	_, _ = d.CreateDeviceMetaTags(device.UUID, body.DeviceMetaTags)
	if device.Name != body.DeviceName || !reflect.DeepEqual(device.Tags, body.DeviceTags) {
		device.Name = body.DeviceName
		device.Tags = body.DeviceTags
		device, err = d.UpdateDevice(device.UUID, device)
		return network, device, err
	}
	return network, device, nil
}
