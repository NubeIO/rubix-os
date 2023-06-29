package main

import (
	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/plugin/pluginapi"
	"github.com/NubeIO/rubix-os/src/cachestore"
	"github.com/NubeIO/rubix-os/src/dbhandler"
)

const path = "inauroazuresync" // must be unique across all plugins
const name = "inauroazuresync" // must be unique across all plugins
const description = "Inauro Azure IoT Hub Sync"
const author = "Nube iO"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "ip"

const pluginType = "database"
const allowConfigWrite = false
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "inauroazuresync"
const transportType = "ip"

// Instance is plugin instance
type Instance struct {
	config                *Config
	enabled               bool
	running               bool
	fault                 bool
	basePath              string
	db                    dbhandler.Handler
	store                 cachestore.Handler
	bus                   eventbus.BusService
	pluginUUID            string
	GatewaySensorTracking InauroGatewaySensorTrackingList // a list of deviceUUIDs.  The length of this is sent in the gateway payload
}

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() pluginapi.Info {
	return pluginapi.Info{
		ModulePath:   path,
		Name:         name,
		Description:  description,
		Author:       author,
		Website:      webSite,
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
