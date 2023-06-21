package modbuschema

import "github.com/NubeIO/rubix-os/schema/schema"

type DeviceSchema struct {
	UUID           schema.UUID                     `json:"uuid"`
	Name           schema.Name                     `json:"name"`
	Description    schema.Description              `json:"description"`
	Enable         schema.Enable                   `json:"enable"`
	TransportType  schema.TransportType            `json:"transport_type"`
	AddressId      schema.AddressId                `json:"address_id"`
	Host           schema.Host                     `json:"host"`
	Port           schema.Port                     `json:"port"`
	FastPollRate   schema.FastPollRate             `json:"fast_poll_rate"`
	NormalPollRate schema.NormalPollRate           `json:"normal_poll_rate"`
	SlowPollRate   schema.SlowPollRate             `json:"slow_poll_rate"`
	ZeroMode       schema.ZeroMode                 `json:"zero_mode"`
	HistoryEnable  schema.HistoryEnableDefaultTrue `json:"history_enable"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	m.Port.Default = 502
	schema.Set(m)
	return m
}
