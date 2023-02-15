package main

import (
	"context"

	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
)

const pluginPath = "lorawan" // must be unique across all plugins
const pluginName = "lorawan" // must be unique across all plugins
const description = "lorawan api"
const author = "Shiny380"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"

const pluginType = "protocol"
const allowConfigWrite = false
const isNetwork = true
const maxAllowedNetworks int = 1
const networkType = "lorawan"
const transportType = "ip" // serial, ip

// Instance is plugin instance
type Instance struct {
	config      *Config
	enabled     bool
	running     bool
	fault       bool
	basePath    string
	db          dbhandler.Handler
	store       cachestore.Handler
	bus         eventbus.BusService
	pluginUUID  string
	networkUUID string
	ctx         context.Context
	cancel      func()
	REST        csrest.RestClient
	csConnected bool
	deviceEUIs  []string
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() pluginapi.Info {
	return pluginapi.Info{
		ModulePath:   pluginPath,
		Name:         pluginName,
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
