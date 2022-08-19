package jsonschema

import (
	"github.com/NubeIO/lib-schema/schema"
)

type NetworkSchema struct {
	UUID           schema.UUID           `json:"uuid"`
	Name           schema.Name           `json:"name"`
	Description    schema.Description    `json:"description"`
	Enable         schema.Enable         `json:"enable"`
	PluginName     schema.PluginName     `json:"plugin_name"`
	TransportType  schema.TransportType  `json:"transport_type"`
	SerialPort     schema.SerialPort     `json:"serial_port"`
	SerialBaudRate schema.SerialBaudRate `json:"serial_baud_rate"`
	SerialParity   schema.SerialParity   `json:"serial_parity"`
	SerialDataBits schema.SerialDataBits `json:"serial_data_bits"`
	SerialStopBits schema.SerialStopBits `json:"serial_stop_bits"`
	SerialTimeout  schema.SerialTimeout  `json:"serial_timeout"`
}

func GetNetworkSchema() *NetworkSchema {
	m := &NetworkSchema{}
	schema.Set(m)
	return m
}
