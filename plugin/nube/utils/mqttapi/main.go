package main

import (
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/mqttclient"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"github.com/NubeDev/flow-framework/src/cachestore"
	"github.com/NubeDev/flow-framework/src/dbhandler"
)

const name = "mqttapi" //must be unique across all plugins
const description = "mqttapi"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "mqtt"

const pluginType = "network"
const allowConfigWrite = false
const isNetwork = false
const maxAllowedNetworks = 1
const networkType = "ip"
const transportType = "ip" //serial, ip

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
	mqtt        *mqttclient.Client
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:   name,
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
