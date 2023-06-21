package loraschema

import "github.com/NubeIO/rubix-os/schema/schema"

const (
	DeviceModelTHLM         = "THLM"
	DeviceModelTHL          = "THL"
	DeviceModelTH           = "TH"
	DeviceModelMicroEdgeV1  = "MicroEdgeV1"
	DeviceModelMicroEdgeV2  = "MicroEdgeV2"
	DeviceModelZiptHydroTap = "ZipHydroTap"
)

type DeviceSchema struct {
	UUID          schema.UUID                     `json:"uuid"`
	Name          schema.Name                     `json:"name"`
	Description   schema.Description              `json:"description"`
	Enable        schema.Enable                   `json:"enable"`
	AddressUUID   schema.AddressUUID              `json:"address_uuid"`
	Model         schema.Model                    `json:"model"`
	HistoryEnable schema.HistoryEnableDefaultTrue `json:"history_enable"`
}

func GetDeviceSchema() *DeviceSchema {
	models := []string{DeviceModelTHLM, DeviceModelTHL, DeviceModelTH, DeviceModelMicroEdgeV1, DeviceModelMicroEdgeV2, DeviceModelZiptHydroTap}
	m := &DeviceSchema{}
	m.AddressUUID.Min = 8
	m.AddressUUID.Max = 8
	m.Model.EnumName = models
	m.Model.Options = models
	schema.Set(m)
	return m
}
