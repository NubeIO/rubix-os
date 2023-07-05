package main

import (
	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/plugin/pluginapi"
	"github.com/NubeIO/rubix-os/src/cachestore"
	"github.com/NubeIO/rubix-os/src/dbhandler"
	"github.com/patrickmn/go-cache"
)

const path = "cpsprocessing" // must be unique across all plugins
const name = "cpsprocessing" // must be unique across all plugins
const description = "cps data processing plugin"
const author = "md"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "processing"
const DefaultExpiration = cache.DefaultExpiration

const pluginType = "project"
const allowConfigWrite = true
const isNetwork = true
const maxAllowedNetworks = 100
const networkType = "cpsprocessing"
const transportType = "ip" // serial, ip

// Instance is plugin instance
type Instance struct {
	config     *Config
	enabled    bool
	running    bool
	fault      bool
	basePath   string
	db         dbhandler.Handler
	store      cachestore.Handler
	bus        eventbus.BusService
	pluginUUID string
	pluginName string
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() pluginapi.Info {
	return pluginapi.Info{
		ModulePath:   path,
		Name:         name,
		Description:  description,
		Author:       author,
		Website:      webSite,
		HasNetwork:   true,
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
