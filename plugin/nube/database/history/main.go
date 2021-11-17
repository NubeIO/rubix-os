package main

import (
	"github.com/NubeIO/flow-framework/eventbus"
	lwrest "github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/restclient"
	"github.com/NubeIO/flow-framework/plugin/plugin-api"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/flow-framework/src/jobs"
	"github.com/patrickmn/go-cache"
)

const path = "history" //must be unique across all plugins
const name = "history" //must be unique across all plugins
const description = "history"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"
const DefaultExpiration = cache.DefaultExpiration

const pluginType = "database"
const allowConfigWrite = false
const isNetwork = false
const maxAllowedNetworks = 0
const networkType = "na"
const transportType = "na"

// Instance is plugin instance
type Instance struct {
	enabled     bool
	basePath    string
	db          dbhandler.Handler
	store       cachestore.Handler
	bus         eventbus.BusService
	pluginUUID  string
	networkUUID string
	REST        *lwrest.RestClient
	jobs        jobs.Jobs
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
func NewFlowPluginInstance() plugin.Plugin {
	return &Instance{}
}

//main will not let main run
func main() {
	panic("this should be built as plugin")
}
