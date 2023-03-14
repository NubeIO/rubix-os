package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type AutoMapping struct {
	FlowNetworkUUID   string                  `json:"flown_network_uuid"`
	StreamUUID        string                  `json:"stream_uuid"`
	ProducerUUID      string                  `json:"product_uuid"`
	NetworkGlobalUUID string                  `json:"network_global_uuid"`
	NetworkName       string                  `json:"network_name"`
	NetworkTags       []*model.Tag            `json:"network_tags"`
	NetworkMetaTags   []*model.NetworkMetaTag `json:"network_meta_tags"`
	DeviceName        string                  `json:"device_name"`
	DeviceTags        []*model.Tag            `json:"device_tags"`
	DeviceMetaTags    []*model.DeviceMetaTag  `jso:"device_meta_tags"`
	PointName         string                  `json:"point_name"`
	PointTags         []*model.Tag            `json:"point_tags"`
	PointMetaTags     []*model.PointMetaTag   `json:"point_meta_tags"`
	IsLocal           bool                    `json:"is_local"`
}
