package main

import (
	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/plugin/pluginapi"
	"github.com/NubeIO/rubix-os/src/cachestore"
	"github.com/NubeIO/rubix-os/src/dbhandler"
)

const path = "maplora" // must be unique across all plugins
const name = "maplora" // must be unique across all plugins
const description = "Maps Rubix LoRa Raw Service to FF"
const author = "md"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"

const pluginType = "protocols"
const allowConfigWrite = false
const isNetwork = false
const maxAllowedNetworks = 0
const networkType = "N/A"
const transportType = "N/A"

// Instance is plugin instance
type Instance struct {
	config     *Config
	enabled    bool
	basePath   string
	db         dbhandler.Handler
	store      cachestore.Handler
	bus        eventbus.BusService
	pluginUUID string
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() pluginapi.Info {
	return pluginapi.Info{
		ModulePath:   path,
		Name:         name,
		Description:  description,
		Author:       author,
		Website:      webSite,
		ProtocolType: protocolType,
	}
}

// NewFlowPluginInstance creates a plugin instance for a user context.
func NewFlowPluginInstance() pluginapi.Plugin {
	return &Instance{}
}

// main will not let main run
func main() {
	panic("this should be built as plugin")
}
