package main

import (
	"context"

	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
)

const pluginPath = "networklinker" // must be unique across all plugins
const pluginName = "networklinker" // must be unique across all plugins
const description = "Plugin linker/Merger/Combiner"
const author = "Shiny380"
const webSite = "https://www.github.com/NubeIO"

const pluginType = "protocol"
const allowConfigWrite = false
const hasNetwork = true
const networkType = "duplicate"
const transportType = ""
const protocolType = "link"

const UI_SEPARATOR = " <-> "
const INTERNAL_SEPARATOR = ":"

// Instance is plugin instance
type Instance struct {
	config     *Config
	enabled    bool
	basePath   string
	db         dbhandler.Handler
	store      cachestore.Handler
	bus        eventbus.BusService
	pluginUUID string
	ctx        context.Context
	cancel     context.CancelFunc
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() pluginapi.Info {
	return pluginapi.Info{
		ModulePath:   pluginPath,
		Name:         pluginName,
		Description:  description,
		Author:       author,
		Website:      webSite,
		HasNetwork:   hasNetwork,
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
