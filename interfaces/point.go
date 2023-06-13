package interfaces

type PointWithParent struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	DeviceUUID  string `json:"device_uuid"`
	DeviceName  string `json:"device_name"`
	NetworkUUID string `json:"network_uuid"`
	NetworkName string `json:"network_name"`
}

type PointHistoryInterval struct {
	UUID            string   `json:"uuid"`
	HistoryInterval *int     `json:"history_interval,omitempty"`
	Timestamp       string   `json:"timestamp,omitempty"`
	PresentValue    *float64 `json:"present_value,omitempty"`
}

type PointForPostgresSync struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	DeviceUUID   string `json:"device_uuid,omitempty"`
	DeviceName   string `json:"device_name,omitempty"`
	NetworkUUID  string `json:"network_uuid"`
	NetworkName  string `json:"network_name"`
	GlobalUUID   string `json:"global_uuid"`
	HostUUID     string `json:"host_uuid"`
	HostName     string `json:"host_name"`
	GroupUUID    string `json:"group_uuid"`
	GroupName    string `json:"group_name"`
	LocationUUID string `json:"location_uuid"`
	LocationName string `json:"location_name"`
}

type PointTagForPostgresSync struct {
	PointUUID string `json:"point_uuid"`
	Tag       string `json:"tag"`
}
