package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type AutoMapping struct {
	FlowNetworkUUID string       `json:"flown_network_uuid"`
	StreamUUID      string       `json:"stream_uuid"`
	ProducerUUID    string       `json:"product_uuid"`
	NetworkUUID     string       `json:"network_uuid"`
	NetworkName     string       `json:"network_name"`
	NetworkTags     []*model.Tag `json:"network_tags"`
	DeviceUUID      string       `json:"device_uuid"`
	DeviceName      string       `json:"device_name"`
	DeviceTags      []*model.Tag `json:"device_tags"`
	PointUUID       string       `json:"point_uuid"`
	PointName       string       `json:"point_name"`
	PointTags       []*model.Tag `json:"point_tags"`
	IsLocal         bool         `json:"is_local"`
}
