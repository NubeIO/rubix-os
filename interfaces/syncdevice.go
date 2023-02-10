package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type SyncDevice struct {
	NetworkUUID     string
	NetworkName     string
	NetworkTags     []*model.Tag
	NetworkMetaTags []*model.NetworkMetaTag
	DeviceUUID      string
	DeviceName      string
	DeviceTags      []*model.Tag
	DeviceMetaTags  []*model.DeviceMetaTag
	FlowNetworkUUID string
	IsLocal         bool
}
