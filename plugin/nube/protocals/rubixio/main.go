package main

import (
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/rubixio/rubixiorest"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
)

const name = "rubix-io" //must be unique across all plugins
const description = "rubix-io api"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"

const pluginType = "protocol"
const allowConfigWrite = true
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "rubix-io"
const transportType = "ip" //serial, ip

// Instance is plugin instance
type Instance struct {
	config         *Config
	enabled        bool
	basePath       string
	db             dbhandler.Handler
	store          cachestore.Handler
	bus            eventbus.BusService
	pluginUUID     string
	networkUUID    string
	rest           *rubixiorest.RestClient
	pollingEnabled bool
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

//main will not let main run
func main() {
	panic("this should be built as plugin")
}
