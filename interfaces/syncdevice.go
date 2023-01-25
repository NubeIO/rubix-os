package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type SyncDevice struct {
	NetworkUUID     string
	NetworkName     string
	NetworkTags     []*model.Tag
	DeviceUUID      string
	DeviceName      string
	DeviceTags      []*model.Tag
	FlowNetworkUUID string
	IsLocal         bool
}
