package main

import (
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/patrickmn/go-cache"
)

const path = "statusmismatchalerts" // must be unique across all plugins
const name = "statusmismatchalerts" // must be unique across all plugins
const description = "statusmismatchalerts plugin"
const author = "md"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "lorawan-modbus"
const DefaultExpiration = cache.DefaultExpiration

const pluginType = "project"
const allowConfigWrite = true
const isNetwork = false
const maxAllowedNetworks = 0
const networkType = "statusmismatchalerts"
const transportType = "serial" // serial, ip

// Instance is plugin instance
type Instance struct {
	config         *Config
	enabled        bool
	basePath       string
	db             dbhandler.Handler
	store          cachestore.Handler
	bus            eventbus.BusService
	pluginUUID     string
	pluginName     string
	networks       []*model.Network
	pollingEnabled bool
	pollingCancel  func()
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