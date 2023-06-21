package rubixioschema

import (
	"github.com/NubeIO/rubix-os/schema/schema"
)

type NetworkSchema struct {
	UUID          schema.UUID          `json:"uuid"`
	Name          schema.Name          `json:"name"`
	Description   schema.Description   `json:"description"`
	Enable        schema.Enable        `json:"enable"`
	PluginName    schema.PluginName    `json:"plugin_name"`
	HistoryEnable schema.HistoryEnable `json:"history_enable"`
}

func GetNetworkSchema() *NetworkSchema {
	m := &NetworkSchema{}
	m.Enable.Default = true
	schema.Set(m)
	return m
}
