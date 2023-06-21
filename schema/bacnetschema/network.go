package bacnetschema

import (
	"github.com/NubeIO/lib-networking/networking"
	"github.com/NubeIO/rubix-os/schema/schema"
)

var nets = networking.New()

type NetworkSchema struct {
	UUID          schema.UUID          `json:"uuid"`
	Name          schema.Name          `json:"name"`
	Description   schema.Description   `json:"description"`
	Enable        schema.Enable        `json:"enable"`
	Interface     schema.Interface     `json:"network_interface"`
	PluginName    schema.PluginName    `json:"plugin_name"`
	HistoryEnable schema.HistoryEnable `json:"history_enable"`
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
