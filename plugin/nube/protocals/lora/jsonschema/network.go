package jsonschema

import (
	"github.com/NubeIO/lib-schema/schema"
)

type NetworkSchema struct {
	UUID                         schema.UUID                         `json:"uuid"`
	Name                         schema.Name                         `json:"name"`
	Description                  schema.Description                  `json:"description"`
	Enable                       schema.Enable                       `json:"enable"`
	PluginName                   schema.PluginName                   `json:"plugin_name"`
	SerialPort                   schema.SerialPort                   `json:"serial_port"`
	SerialBaudRate               schema.SerialBaudRate               `json:"serial_baud_rate"`
	AutoMappingNetworksSelection schema.AutoMappingNetworksSelection `json:"auto_mapping_networks_selection"`
	AutoMappingFlowNetworkName   schema.AutoMappingFlowNetworkName   `json:"auto_mapping_flow_network_name"`
	AutoMappingFlowNetworkUUID   schema.AutoMappingFlowNetworkUUID   `json:"auto_mapping_flow_network_uuid"`
	AutoMappingEnableHistories   schema.AutoMappingEnableHistories   `json:"auto_mapping_enable_histories"`
}

func GetNetworkSchema() *NetworkSchema {
	m := &NetworkSchema{}
	m.SerialBaudRate.ReadOnly = true
	m.SerialBaudRate.EnumName = []string{"9600", "38400", "57600", "115200"}
	m.SerialBaudRate.Options = []string{"9600", "38400", "57600", "115200"}
	//Options  []string `json:"options" default:"[\"/dev/ttyAMA0\",\"/data/socat/loRa1\",\"/dev/ttyUSB0\",\"/dev/ttyUSB1\",\"/dev/ttyUSB2\",\"/dev/ttyUSB3\",\"/dev/ttyUSB4\"]"`
	//ports, err := serial.GetPortsList()
	//if err != nil {
	//	log.Errorf("lora: get serial ports for schema err:%s", err.Error())
	//} else {
	//	m.SerialPort.Options = ports
	//	m.SerialPort.EnumName = ports
	//}
	schema.Set(m)
	return m
}
