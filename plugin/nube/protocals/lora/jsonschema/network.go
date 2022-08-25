package jsonschema

import (
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

type NetworkSchema struct {
	UUID                         schema.UUID                         `json:"uuid"`
	Name                         schema.Name                         `json:"name"`
	Description                  schema.Description                  `json:"description"`
	Enable                       schema.Enable                       `json:"enable"`
	PluginName                   schema.PluginName                   `json:"plugin_name"`
	SerialPort                   schema.SerialPort                   `json:"serial_port"`
	AutoMappingNetworksSelection schema.AutoMappingNetworksSelection `json:"auto_mapping_networks_selection"`
	AutoMappingFlowNetworkName   schema.AutoMappingFlowNetworkName   `json:"auto_mapping_flow_network_name"`
	AutoMappingFlowNetworkUUID   schema.AutoMappingFlowNetworkUUID   `json:"auto_mapping_flow_network_uuid"`
	AutoMappingEnableHistories   schema.AutoMappingEnableHistories   `json:"auto_mapping_enable_histories"`
}

func GetNetworkSchema() *NetworkSchema {
	m := &NetworkSchema{}
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Errorf("loraraw: get serial ports for schema err:%s", err.Error())
	} else {
		m.SerialPort.Options = ports
		m.SerialPort.EnumName = ports

	}
	schema.Set(m)
	return m
}
