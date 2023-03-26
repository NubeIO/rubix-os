package interfaces

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

type AutoMappingNetwork struct {
	GlobalUUID      string                  `json:"global_uuid"`
	UUID            string                  `json:"uuid"`
	Name            string                  `json:"name"`
	Tags            []*model.Tag            `json:"tags"`
	MetaTags        []*model.NetworkMetaTag `json:"meta_tags"`
	Devices         []*AutoMappingDevice    `json:"devices"`
	FlowNetworkUUID string                  `json:"flown_network_uuid"`
}

type AutoMappingNetworkError struct {
	UUID    string                    `json:"uuid"`
	Error   *string                   `json:"error"`
	Devices []*AutoMappingDeviceError `json:"devices"`
}

type AutoMappingDevice struct {
	UUID       string                 `json:"uuid"`
	Name       string                 `json:"name"`
	Tags       []*model.Tag           `json:"tags"`
	MetaTags   []*model.DeviceMetaTag `json:"meta_tags"`
	Points     []*AutoMappingPoint    `json:"points"`
	StreamUUID string                 `json:"stream_uuid"`
}

type AutoMappingDeviceError struct {
	Name      string                      `json:"device_name"`
	Error     *string                     `json:"error"`
	Points    []*AutoMappingPointError    `json:"points"`
	Consumers []*AutoMappingConsumerError `json:"consumers"`
	Writers   []*AutoMappingWriterError   `json:"writers"`
}

type AutoMappingPoint struct {
	UUID         string                `json:"uuid"`
	Name         string                `json:"name"`
	Tags         []*model.Tag          `json:"tags"`
	MetaTags     []*model.PointMetaTag `json:"meta_tags"`
	ProducerUUID string                `json:"product_uuid"`
}

type AutoMappingPointError struct {
	Name  string  `json:"name"`
	Error *string `json:"error"`
}

type AutoMappingConsumerError struct {
	Name  string  `json:"name"`
	Error *string `json:"error"`
}

type AutoMappingWriterError struct {
	Name  string  `json:"name"`
	Error *string `json:"error"`
}
