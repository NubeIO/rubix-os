package systemschema

import "github.com/NubeIO/rubix-os/schema/schema"

type DeviceSchema struct {
	UUID          schema.UUID                     `json:"uuid"`
	Name          schema.Name                     `json:"name"`
	Description   schema.Description              `json:"description"`
	Enable        schema.Enable                   `json:"enable"`
	HistoryEnable schema.HistoryEnableDefaultTrue `json:"history_enable"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	schema.Set(m)
	return m
}
