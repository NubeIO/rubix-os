package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type Level int

const (
	Network Level = iota
	Device
	Point
)

func (s Level) String() string {
	switch s {
	case Network:
		return "Network"
	case Device:
		return "Device"
	}
	return "Point"
}

type AutoMappingNetwork struct {
	GlobalUUID      string                  `json:"global_uuid"`
	UUID            string                  `json:"uuid"`
	Name            string                  `json:"name"`
	Tags            []*model.Tag            `json:"tags"`
	MetaTags        []*model.NetworkMetaTag `json:"meta_tags"`
	Devices         []*AutoMappingDevice    `json:"devices"`
	FlowNetworkUUID string                  `json:"flown_network_uuid"`
}

type AutoMappingError struct {
	NetworkUUID string `json:"network_uuid"`
	DeviceUUID  string `json:"device_uuid"`
	PointUUID   string `json:"point_uuid"`
	HasError    bool   `json:"has_error"`
	Error       string `json:"error"`
	Level       Level  `json:"level"`
}

type AutoMappingDevice struct {
	UUID       string                 `json:"uuid"`
	Name       string                 `json:"name"`
	Tags       []*model.Tag           `json:"tags"`
	MetaTags   []*model.DeviceMetaTag `json:"meta_tags"`
	Points     []*AutoMappingPoint    `json:"points"`
	StreamUUID string                 `json:"stream_uuid"`
}

type AutoMappingPoint struct {
	UUID         string                `json:"uuid"`
	Name         string                `json:"name"`
	Tags         []*model.Tag          `json:"tags"`
	MetaTags     []*model.PointMetaTag `json:"meta_tags"`
	ProducerUUID string                `json:"product_uuid"`
}
