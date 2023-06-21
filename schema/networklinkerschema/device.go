package networklinkerschema

import "github.com/NubeIO/rubix-os/schema/schema"

type DeviceSchema struct {
	Enable      schema.Enable `json:"enable" deafult:"true"`
	Name        schema.Name   `json:"name"`
	AddressUUID struct {
		Type    string               `json:"type" default:"string"`
		Title   string               `json:"title" default:"address_uuid"`
		Options []schema.OptionOneOf `json:"oneOf"`
		Help    string               `json:"help" default:"address_uuid"`
	} `json:"address_uuid"`
	HistoryEnable schema.HistoryEnableDefaultTrue `json:"history_enable"`
}

func GetDeviceSchema() *DeviceSchema {
	m := &DeviceSchema{}
	schema.Set(m)
	return m
}
