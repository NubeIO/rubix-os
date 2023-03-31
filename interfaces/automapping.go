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

type AutoMappingResponse struct {
	NetworkUUID string        `json:"network_uuid"`
	DeviceUUID  string        `json:"device_uuid"`
	PointUUID   string        `json:"point_uuid"`
	HasError    bool          `json:"has_error"`
	Error       string        `json:"error"`
	Level       Level         `json:"level"`
	SyncWriters []*SyncWriter `json:"sync_writers"`
}

type SyncWriter struct {
	ProducerUUID      string
	WriterUUID        string
	FlowFrameworkUUID string
	PointUUID         string
	PointName         string
}

type AutoMapping struct {
	GlobalUUID      string                `json:"global_uuid"`
	FlowNetworkUUID string                `json:"flow_network_uuid"`
	Level           Level                 `json:"level"`
	Networks        []*AutoMappingNetwork `json:"networks"`
}

type AutoMappingNetwork struct {
	Enable   bool                    `json:"enable"`
	UUID     string                  `json:"uuid"`
	Name     string                  `json:"name"`
	Tags     []*model.Tag            `json:"tags"`
	MetaTags []*model.NetworkMetaTag `json:"meta_tags"`
	Devices  []*AutoMappingDevice    `json:"devices"`
}

type AutoMappingDevice struct {
	Enable          bool                   `json:"enable"`
	UUID            string                 `json:"uuid"`
	Name            string                 `json:"name"`
	Tags            []*model.Tag           `json:"tags"`
	MetaTags        []*model.DeviceMetaTag `json:"meta_tags"`
	Points          []*AutoMappingPoint    `json:"points"`
	StreamUUID      string                 `json:"stream_uuid"`
	StreamCloneUUID string                 `json:"stream_clone_uuid"`
}

type AutoMappingPoint struct {
	Enable       bool                  `json:"enable"`
	UUID         string                `json:"uuid"`
	Name         string                `json:"name"`
	Tags         []*model.Tag          `json:"tags"`
	MetaTags     []*model.PointMetaTag `json:"meta_tags"`
	ProducerUUID string                `json:"product_uuid"`
}
