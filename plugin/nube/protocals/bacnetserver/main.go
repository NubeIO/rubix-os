package main

import (
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"github.com/NubeDev/flow-framework/src/cachestore"
	"github.com/NubeDev/flow-framework/src/dbhandler"
)

const path = "bacnetserver" //must be unique across all plugins
const name = "bacnetserver" //must be unique across all plugins
const description = "bacnet server api to nube bacnet stack"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"

const pluginType = "protocol"
const allowConfigWrite = false
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "bacnet"
const transportType = "ip" //serial, ip
const ip = "0.0.0.0"
const port = "1717"

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
