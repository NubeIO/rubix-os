package masterschema

import (
	"github.com/NubeIO/lib-networking/networking"
	"github.com/NubeIO/rubix-os/schema/schema"
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
	UUID          schema.UUID          `json:"uuid"`
	Name          schema.Name          `json:"name"`
	Description   schema.Description   `json:"description"`
	Enable        schema.Enable        `json:"enable"`
	Port          schema.Port          `json:"port"`
	Interface     schema.Interface     `json:"network_interface"`
	PluginName    schema.PluginName    `json:"plugin_name"`
	MaxPollRate   schema.MaxPollRate   `json:"max_poll_rate"`
	HistoryEnable schema.HistoryEnable `json:"history_enable"`
}

func GetNetworkSchema() *NetworkSchema {
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
	schema.Set(m)
	return m
}
