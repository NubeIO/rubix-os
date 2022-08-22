package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/services/pollqueue"
	"github.com/NubeIO/flow-framework/utils/float"
)

func (inst *Instance) Enable() error {
	fmt.Println("11111")
	inst.bacnetDebugMsg("Polling Enable()")
	inst.enabled = true
	inst.pluginName = name
	inst.setUUID()

	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	if nets != nil {
		inst.networks = nets
	} else if err != nil {
		inst.networks = nil
	}
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
			err := inst.BACnetMasterPolling()
			if err != nil {
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
	return nil
}
