package interfaces

type PublishPointList struct {
	PluginPath  string `json:"plugin_path"`
	NetworkName string `json:"network_name"`
	DeviceName  string `json:"device_name"`
	PointUUID   string `json:"point_uuid"`
	PointName   string `json:"point_name"`
}
