package loraschema

import (
	"github.com/NubeIO/rubix-os/schema/schema"
)

type NetworkSchema struct {
	UUID           schema.UUID           `json:"uuid"`
	Name           schema.Name           `json:"name"`
	Description    schema.Description    `json:"description"`
	Enable         schema.Enable         `json:"enable"`
	PluginName     schema.PluginName     `json:"plugin_name"`
	SerialPort     schema.SerialPortLora `json:"serial_port"`
	SerialBaudRate schema.SerialBaudRate `json:"serial_baud_rate"`
	HistoryEnable  schema.HistoryEnable  `json:"history_enable"`
}

func GetNetworkSchema() *NetworkSchema {
	m := &NetworkSchema{}
	schema.Set(m)
	return m
}
