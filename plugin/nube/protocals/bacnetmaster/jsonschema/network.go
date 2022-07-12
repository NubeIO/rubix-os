package jsonschema

import (
	"github.com/NubeIO/lib-networking/networking"
	"github.com/NubeIO/lib-schema/schema"
)

var nets = networking.New()

type NetworkSchema struct {
	UUID                         schema.UUID                         `json:"uuid"`
	Name                         schema.Name                         `json:"name"`
	Description                  schema.Description                  `json:"description"`
	Enable                       schema.Enable                       `json:"enable"`
	Port                         schema.Port                         `json:"port"`
	Interface                    schema.Interface                    `json:"network_interface"`
	PluginName                   schema.PluginName                   `json:"plugin_name"`
	AutoMappingNetworksSelection schema.AutoMappingNetworksSelection `json:"auto_mapping_networks_selection"`
	AutoMappingFlowNetworkName   schema.AutoMappingFlowNetworkName   `json:"auto_mapping_flow_network_name"`
	AutoMappingFlowNetworkUUID   schema.AutoMappingFlowNetworkUUID   `json:"auto_mapping_flow_network_uuid"`
	AutoMappingEnableHistories   schema.AutoMappingEnableHistories   `json:"auto_mapping_enable_histories"`
}

func GetNetworkSchema() *NetworkSchema {
	m := &NetworkSchema{}
	schema.Set(m)
	names, err := nets.GetInterfacesNames()
	if err != nil {
		return m
	}
	var out []string
	out = append(out, "eth0")
	out = append(out, "eth1")
	for _, name := range names.Names {
		if name != "lo" {
			out = append(out, name)
		}
	}
	m.Interface.Options = out
	return m
}
