package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"math"
	"math/rand"
	"time"

	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
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

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = float64(int(math.Round(min + rand.Float64()*(max-min))))
	}
	return res
}

// run LoRa plugin loop
func (inst *Instance) run() {
	points, err := inst.db.GetPoints(api.Args{})
	fmt.Println("points", points, err)
	for {
		for _, point := range points {
			val := randFloats(1, 2, 1)[0]
			go func() {
				point.Priority = &model.Priority{P16: &val}
				_, _ = inst.db.UpdatePoint(point.UUID, point, true, false)
				//p, err := inst.db.UpdatePoint(point.UUID, point, true, false)
				//fmt.Println("p", p)
				//fmt.Println("err", err)
			}()
			time.Sleep(1 * time.Millisecond)
		}
		time.Sleep(1 * time.Millisecond)
	}
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
