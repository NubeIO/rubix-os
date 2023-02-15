package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/services/pollqueue"
	"github.com/NubeIO/flow-framework/utils/float"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) Enable() error {
	inst.bacnetDebugMsg("Polling Enable()")
	inst.enabled = true
	inst.fault = false
	inst.running = false
	inst.pluginName = name
	inst.setUUID()

	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	if err != nil {
		inst.fault = true
		inst.bacnetErrorMsg("enable plugin get networks: %v\n", err)
	}
	log.Infof("bacnet-master: enable plugin networks count: %d pluginUUID: %s", len(nets), inst.pluginUUID)
	inst.initBacStore()

	if inst.config.EnablePolling {
		if !inst.pollingEnabled {
			inst.pollingEnabled = true
			inst.NetworkPollManagers = make([]*pollqueue.NetworkPollManager, 0) // This will delete any existing NetworkPollManagers (if enable is called multiple times, it will rebuild the queues).
			for _, net := range nets {                                          // Create a new Poll Manager for each network in the plugin.
				conf := inst.GetConfig().(*Config)
				pollQueueConfig := pollqueue.Config{EnablePolling: conf.EnablePolling, LogLevel: conf.LogLevel}
				pollManager := NewPollManager(&pollQueueConfig, &inst.db, net.UUID, inst.pluginUUID, inst.pluginName, float.NonNil(net.MaxPollRate))
				// inst.modbusDebugMsg("net")
				// inst.modbusDebugMsg("%+v\n", net)
				// inst.modbusDebugMsg("pollManager")
				// inst.modbusDebugMsg("%+v\n", pollManager)
				pollManager.StartPolling()
				inst.NetworkPollManagers = append(inst.NetworkPollManagers, pollManager)
			}

			// TODO: VERIFY POLLING WITHOUT GO ROUTINE WRAPPER
			inst.running = true
			err := inst.BACnetMasterPolling()
			if err != nil {
				inst.running = false
				inst.fault = true
				inst.bacnetErrorMsg("POLLING ERROR on routine: %v\n", err)
			}
		}
	}
	return nil
}

func (inst *Instance) Disable() error {
	inst.bacnetDebugMsg("Polling Disable()")
	inst.enabled = false
	if inst.pollingEnabled {
		inst.pollingEnabled = false
		inst.pollingCancel()
		inst.pollingCancel = nil
		for _, pollMan := range inst.NetworkPollManagers {
			pollMan.StopPolling()
			// inst.closeBacnetStoreNetwork(pollMan.FFNetworkUUID)  // TODO: this causes FF to lock up
		}
		inst.NetworkPollManagers = make([]*pollqueue.NetworkPollManager, 0)
	}
	inst.running = false
	inst.fault = false
	return nil
}
