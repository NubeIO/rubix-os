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
