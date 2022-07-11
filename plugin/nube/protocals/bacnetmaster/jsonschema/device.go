package jsonschema

import "github.com/NubeIO/lib-schema/schema"

type DeviceSchema struct {
	UUID           schema.UUID           `json:"uuid"`
	Name           schema.Name           `json:"name"`
	Description    schema.Description    `json:"description"`
	Enable         schema.Enable         `json:"enable"`
	Host           schema.Host           `json:"host"`
	Port           schema.Port           `json:"port"`
	DeviceObjectId schema.DeviceObjectId `json:"device_object_id"`
	NetworkNumber  schema.NetworkNumber  `json:"network_number"`
	DeviceMac      schema.DeviceMac      `json:"device_mac"`
	MaxADPU        schema.MaxADPU        `json:"max_adpu"`
	Segmentation   schema.Segmentation   `json:"segmentation"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	schema.Set(m)
	return m
}
