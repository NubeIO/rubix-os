package jsonschema

import "github.com/NubeIO/lib-schema/schema"

type DeviceSchema struct {
	UUID        schema.UUID        `json:"uuid"`
	Name        schema.Name        `json:"name"`
	Description schema.Description `json:"description"`
	Enable      schema.Enable      `json:"enable"`
	AddressUUID schema.AddressUUID `json:"address_uuid"`
	Model       schema.Model       `json:"model"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	m.AddressUUID.Default = "88B2FC1A"
	m.AddressUUID.Min = 8
	m.AddressUUID.Max = 8
	m.Model.Default = "THLM"
	m.Model.EnumName = []string{"THLM", "THL", "TH", "TL", "MicroEdge", "ZipHydroTap"}
	m.Model.Options = []string{"THLM", "THL", "TH", "TL", "MicroEdge", "ZipHydroTap"}
	schema.Set(m)
	return m
}
