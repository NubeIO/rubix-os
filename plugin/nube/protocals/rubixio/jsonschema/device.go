package jsonschema

import "github.com/NubeIO/lib-schema/schema"

type DeviceSchema struct {
	UUID        schema.UUID        `json:"uuid"`
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
	Enable      schema.Enable      `json:"enable"`
	Host        schema.Host        `json:"host"`
	Port        schema.Port        `json:"port"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	m.Host.Default = "0.0.0.0"
	m.Port.Default = 5001
	m.Enable.Default = true
	schema.Set(m)
	return m
}
