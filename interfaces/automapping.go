package interfaces

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

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

type AutoMappingScheduleResponse struct {
	ScheduleUUID string        `json:"schedule_uuid"`
	HasError     bool          `json:"has_error"`
	Error        string        `json:"error"`
	SyncWriters  []*SyncWriter `json:"sync_writers"`
}

type SyncWriter struct {
	ProducerUUID      string
	WriterUUID        string
	FlowFrameworkUUID string
	UUID              string
	Name              string
}

type AutoMapping struct {
	GlobalUUID      string                 `json:"global_uuid"`
	FlowNetworkUUID string                 `json:"flow_network_uuid"`
	Level           Level                  `json:"level"`
	Networks        []*AutoMappingNetwork  `json:"networks"`
	Schedules       []*AutoMappingSchedule `json:"schedules"`
}

type AutoMappingNetwork struct {
	Enable            bool                    `json:"enable"`
	AutoMappingEnable bool                    `json:"auto_mapping_enable"`
	UUID              string                  `json:"uuid"`
	Name              string                  `json:"name"`
	Tags              []*model.Tag            `json:"tags"`
	MetaTags          []*model.NetworkMetaTag `json:"meta_tags"`
	Devices           []*AutoMappingDevice    `json:"devices"`
	CreateNetwork     bool                    `json:"create_network"`
}

type AutoMappingDevice struct {
	Enable            bool                   `json:"enable"`
	AutoMappingEnable bool                   `json:"auto_mapping_enable"`
	UUID              string                 `json:"uuid"`
	Name              string                 `json:"name"`
	Tags              []*model.Tag           `json:"tags"`
	MetaTags          []*model.DeviceMetaTag `json:"meta_tags"`
	Points            []*AutoMappingPoint    `json:"points"`
	StreamUUID        string                 `json:"stream_uuid"`
}

type AutoMappingPoint struct {
	Enable            bool                  `json:"enable"`
	AutoMappingEnable bool                  `json:"auto_mapping_enable"`
	EnableWriteable   bool                  `json:"enable_writeable"`
	UUID              string                `json:"uuid"`
	Name              string                `json:"name"`
	Tags              []*model.Tag          `json:"tags"`
	MetaTags          []*model.PointMetaTag `json:"meta_tags"`
	ProducerUUID      string                `json:"product_uuid"`
	Priority          model.Priority        `json:"priority"`
}

type AutoMappingSchedule struct {
	Enable            bool   `json:"enable"`
	AutoMappingEnable bool   `json:"auto_mapping_enable"`
	UUID              string `json:"uuid"`
	Name              string `json:"name"`
	StreamUUID        string `json:"stream_uuid"`
	ProducerUUID      string `json:"product_uuid"`
	CreateSchedule    bool   `json:"create_schedule"`
}
