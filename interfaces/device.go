package interfaces

type DeviceTagForPostgresSync struct {
	DeviceUUID string `json:"device_uuid"`
	Tag        string `json:"tag"`
}
