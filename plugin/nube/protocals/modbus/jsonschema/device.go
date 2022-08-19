package jsonschema

import "github.com/NubeIO/lib-schema/schema"

type DeviceSchema struct {
	UUID           schema.UUID           `json:"uuid"`
	Name           schema.Name           `json:"name"`
	Description    schema.Description    `json:"description"`
	Enable         schema.Enable         `json:"enable"`
	TransportType  schema.TransportType  `json:"transport_type"`
	AddressId      schema.AddressId      `json:"address_id"`
	Host           schema.Host           `json:"host"`
	Port           schema.Port           `json:"port"`
	FastPollRate   schema.FastPollRate   `json:"fast_poll_rate"`
	NormalPollRate schema.NormalPollRate `json:"normal_poll_rate"`
	SlowPollRate   schema.SlowPollRate   `json:"slow_poll_rate"`
	ZeroMode       schema.ZeroMode       `json:"zero_mode"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	schema.Set(m)
	return m
}
