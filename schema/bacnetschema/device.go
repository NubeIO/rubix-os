package bacnetschema

import "github.com/NubeIO/rubix-os/schema/schema"

type DeviceSchema struct {
	UUID           schema.UUID                     `json:"uuid"`
	Name           schema.Name                     `json:"name"`
	Description    schema.Description              `json:"description"`
	Enable         schema.Enable                   `json:"enable"`
	DeviceObjectId schema.DeviceObjectId           `json:"device_object_id"`
	FastPollRate   schema.FastPollRate             `json:"fast_poll_rate"`
	NormalPollRate schema.NormalPollRate           `json:"normal_poll_rate"`
	SlowPollRate   schema.SlowPollRate             `json:"slow_poll_rate"`
	DeviceTimeout  schema.TimeoutSecs              `json:"device_timeout"`
	HistoryEnable  schema.HistoryEnableDefaultTrue `json:"history_enable"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	schema.Set(m)
	return m
}
