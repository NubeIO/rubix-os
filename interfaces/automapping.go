package interfaces

type AutoMapping struct {
	FlowNetworkUUID string `json:"flown_network_uuid"`
	StreamUUID      string `json:"stream_uuid"`
	ProducerUUID    string `json:"product_uuid"`
	NetworkUUID     string `json:"network_uuid"`
	NetworkName     string `json:"network_name"`
	DeviceUUID      string `json:"device_uuid"`
	DeviceName      string `json:"device_name"`
	PointUUID       string `json:"point_uuid"`
	PointName       string `json:"point_name"`
	IsLocal         bool   `json:"is_local"`
}
