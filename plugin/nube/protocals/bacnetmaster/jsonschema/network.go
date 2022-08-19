package jsonschema

import (
	"github.com/NubeIO/lib-networking/networking"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

var nets = networking.New()

type AutoMappingFlowNetworkUUID struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"select flow network for mapping"`
	Options  []string `json:"enum" default:"[]"`
	EnumName []string `json:"enumNames" default:"[]"`
}

type AutoMappingNetworksSelection struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:"enable mapping"`
	Options  []string `json:"enum" default:"[\"disable\",\"self-mapping\",\"bacnet\"]"`
	EnumName []string `json:"enumNames" default:"[\"disable\",\"self-mapping\",\"bacnet\"]"`
	Default  string   `json:"default" default:"disable"`
}

type NetworkSchema struct {
	UUID           schema.UUID           `json:"uuid"`
	Name           schema.Name           `json:"name"`
	Description    schema.Description    `json:"description"`
	Enable         schema.Enable         `json:"enable"`
	Port           schema.Port           `json:"port"`
	Interface      schema.Interface      `json:"network_interface"`
	PluginName     schema.PluginName     `json:"plugin_name"`
	FastPollRate   schema.FastPollRate   `json:"fast_poll_rate"`
	NormalPollRate schema.NormalPollRate `json:"normal_poll_rate"`
	SlowPollRate   schema.SlowPollRate   `json:"slow_poll_rate"`

	//AutoMappingNetworksSelection schema.AutoMappingNetworksSelection `json:"auto_mapping_networks_selection"`
	//AutoMappingFlowNetworkUUID   AutoMappingFlowNetworkUUID          `json:"auto_mapping_flow_network_uuid"`
}

func GetNetworkSchema(flows []*model.FlowNetwork) *NetworkSchema {
	m := &NetworkSchema{}
	m.Port.Default = 47808
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
	//var flowNames []string
	//var flowUUIDS []string
	//for _, flow := range flows {
	//	flowNames = append(flowNames, flow.Name)
	//	flowUUIDS = append(flowUUIDS, flow.UUID)
	//}
	//

	//m.AutoMappingFlowNetworkUUID.EnumName = flowNames
	//m.AutoMappingFlowNetworkUUID.Options = flowUUIDS

	schema.Set(m)
	return m
}
