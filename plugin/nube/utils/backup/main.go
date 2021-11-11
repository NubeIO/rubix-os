package main

import (
	"github.com/NubeDev/flow-framework/eventbus"
	min "github.com/NubeDev/flow-framework/plugin/nube/utils/backup/minio"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"github.com/NubeDev/flow-framework/src/cachestore"
	"github.com/NubeDev/flow-framework/src/dbhandler"
	"github.com/patrickmn/go-cache"
)

const name = "backup" //must be unique across all plugins
const description = "backup"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "git"
const DefaultExpiration = cache.DefaultExpiration

const pluginType = "backup"
const allowConfigWrite = false
const isNetwork = false
const maxAllowedNetworks = 1
const networkType = ""
const transportType = "" //serial, ip

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
	minioClient min.MinioClient
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
func NewFlowPluginInstance() plugin.Plugin {
	return &Instance{}

}

//main will not let main run
func main() {
	panic("this should be built as plugin")
}
