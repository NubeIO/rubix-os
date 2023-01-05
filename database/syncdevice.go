package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) SyncDevice(body *interfaces.SyncDevice) (*model.Device, error) {
	syncNetwork := &model.SyncNetwork{NetworkUUID: body.NetworkUUID, NetworkName: body.NetworkName,
		FlowNetworkUUID: body.FlowNetworkUUID, IsLocal: body.IsLocal}
	network, err := d.SyncNetwork(syncNetwork)
	if err != nil {
		return nil, err
	}
	device, err := d.GetOneDeviceByArgs(api.Args{AutoMappingUUID: nils.NewString(body.DeviceUUID)})
	if err != nil {
		fnc, err := d.GetOneFlowNetworkCloneByArgs(api.Args{SourceUUID: nils.NewString(body.FlowNetworkUUID)})
		if err != nil {
			return nil, err
		}
		deviceModel := &model.Device{}
		deviceModel.Name = body.DeviceName
		deviceModel.Enable = boolean.NewTrue()
		deviceModel.NetworkUUID = network.UUID
		deviceModel.AutoMappingUUID = body.DeviceUUID
		deviceModel.AutoMappingFlowNetworkUUID = fnc.UUID
		deviceModel.AutoMappingFlowNetworkName = fnc.Name
		return d.CreateDevice(deviceModel)
	}
	if device.Name != body.DeviceName {
		device.Name = body.DeviceName
		return d.UpdateDevice(device.UUID, device, false)
	}
	return device, nil
}