package main

import (
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/config"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/edgerest"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

const name = "edge28" // must be unique across all plugins
const description = "edge28 api"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"

const pluginType = "protocol"
const allowConfigWrite = true
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "edge28"
const transportType = "ip" // serial, ip

// Instance is plugin instance
type Instance struct {
	config         *config.Config
	enabled        bool
	running        bool
	fault          bool
	basePath       string
	db             dbhandler.Handler
	store          cachestore.Handler
	bus            eventbus.BusService
	pluginUUID     string
	pluginName     string
	networks       []*model.Network
	pollingEnabled bool
	pollingCancel  func()
	rest           *edgerest.RestClient
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() pluginapi.Info {
	return pluginapi.Info{
		ModulePath:   name,
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
