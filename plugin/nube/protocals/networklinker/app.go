package main

import (
	"context"
	"github.com/NubeIO/rubix-os/args"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (inst *Instance) syncPointsLoopTEMPORARY(ctx context.Context) {
	inst.syncPointValues()
	inst.syncNetworkDevicePoints()
	for {
		select {
		case <-ctx.Done():
			log.Info("networklinker: exiting sync service")
			return
		case <-time.After(time.Duration(inst.config.ValueSyncIntervalSeconds) * time.Second):
			inst.syncPointValues()
		case <-time.After(time.Duration(inst.config.LinkSyncIntervalSeconds) * time.Second):
			inst.syncNetworkDevicePoints()
		}
	}
}

func (inst *Instance) syncPointValues() {
	networks, _ := inst.db.GetNetworksByPlugin(inst.pluginUUID, args.Args{WithDevices: true, WithPoints: true})
	index := 0
	for _, net := range networks {
		netUUIDs := strings.Split(net.AddressUUID, INTERNAL_SEPARATOR)
		network1, _ := inst.db.GetNetwork(netUUIDs[0], args.Args{})
		if !inst.networkIsWriter(network1) {
			index = 0
		} else {
			index = 1
		}
		for _, dev := range net.Devices {
			for _, point := range dev.Points {
				pointUUIDs := strings.Split(*point.AddressUUID, INTERNAL_SEPARATOR)
				origUUID := &pointUUIDs[0]
				if len(pointUUIDs) == 2 && index == 1 {
					origUUID = &pointUUIDs[1]
				}
				inst.syncPointSelected(point, *origUUID)
			}
		}
	}
}

func (inst *Instance) syncNetworkDevicePoints() {
	networks, _ := inst.db.GetNetworksByPlugin(inst.pluginUUID, args.Args{WithDevices: true, WithPoints: true})
	for _, net := range networks {
		for _, dev := range net.Devices {
			inst.syncDevicePoints(dev, nil, nil, nil, nil)
		}
	}
}
