package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type SyncNetwork struct {
	NetworkGlobalUUID string
	NetworkName       string
	NetworkTags       []*model.Tag
	NetworkMetaTags   []*model.NetworkMetaTag
	FlowNetworkUUID   string
	IsLocal           bool
}
