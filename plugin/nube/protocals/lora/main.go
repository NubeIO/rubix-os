package main

import (
	"time"

	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/plugin/pluginapi"
	"github.com/NubeIO/rubix-os/src/cachestore"
	"github.com/NubeIO/rubix-os/src/dbhandler"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

const path = "lora" // must be unique across all plugins
const name = "lora" // must be unique across all plugins
const description = "LoRaRAW"
const author = "ap, dm"
const webSite = "https://www.github.com/NubeIO"
const protocolType = "serial"
const DefaultExpiration = cache.DefaultExpiration

const pluginType = "protocol"
const allowConfigWrite = false
const isNetwork = true
const maxAllowedNetworks = 1
const networkType = "lora"
const transportType = "serial" // serial, ip

// Instance is plugin instance
type Instance struct {
	config        *Config
	enabled       bool
	running       bool
	fault         bool
	basePath      string
	db            dbhandler.Handler
	store         cachestore.Handler
	bus           eventbus.BusService
	pluginUUID    string
	networkUUID   string
	interruptChan chan struct{}
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

// run LoRa plugin loop
func (inst *Instance) run() {
	defer inst.SerialClose()

	for {
		sc, err := inst.SerialOpen()
		if err != nil {
			log.Error("loraraw: error opening serial ", err)
			time.Sleep(5 * time.Second)
			continue
		}
		serialPayloadChan := make(chan string, 1)
		serialCloseChan := make(chan error, 1)
		go sc.Loop(serialPayloadChan, serialCloseChan)

	readLoop:
		for {
			select {
			case <-inst.interruptChan:
				log.Info("loraraw: interrupt received on run")
				return
			case err := <-serialCloseChan:
				log.Error("loraraw: serial connection error: ", err)
				log.Info("loraraw: serial connection attempting to reconnect...")
				break readLoop
			case data := <-serialPayloadChan:
				inst.handleSerialPayload(data)
			}
		}
	}
}
