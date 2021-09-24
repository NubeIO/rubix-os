package main

import (
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"github.com/NubeDev/flow-framework/src/cachestore"
	"github.com/NubeDev/flow-framework/src/dbhandler"
	"github.com/NubeDev/flow-framework/src/eventbus"
	"github.com/patrickmn/go-cache"
)

const path = "lora" //must be unique across all plugins
const name = "lora" //must be unique across all plugins
const description = "lora raw"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "serial"
const DefaultExpiration = cache.DefaultExpiration

const pluginType = "protocol"
const allowConfigWrite = false
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "lora"
const transportType = "serial" //serial, ip

// Instance is plugin instance
type Instance struct {
	config      *Config
	enabled     bool
	basePath    string
	db          dbhandler.Handler
	store       cachestore.Handler
	bus         eventbus.BusService
	pluginUUID  string
	networkUUID string
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:   path,
		Name:         name,
		Description:  description,
		Author:       author,
		Website:      webSite,
		ProtocolType: protocolType,
	}
}

// NewFlowPluginInstance creates a plugin instance for a user context.
func NewFlowPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &Instance{}

}

//main will not let main run
func main() {
	panic("this should be built as plugin")
}
