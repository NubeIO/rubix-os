package interfaces

type MqttPoint struct {
	NetworkName string `json:"network_name"`
	DeviceName  string `json:"device_name"`
	PointName   string `json:"point_name"`
}
