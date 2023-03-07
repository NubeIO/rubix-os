package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type SyncDevice struct {
	NetworkGlobalUUID string
	NetworkName       string
	NetworkTags       []*model.Tag
	NetworkMetaTags   []*model.NetworkMetaTag
	DeviceName        string
	DeviceTags        []*model.Tag
	DeviceMetaTags    []*model.DeviceMetaTag
	FlowNetworkUUID   string
	IsLocal           bool
}
