package main

import (
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/services/pollqueue"
	"github.com/NubeIO/rubix-os/utils/float"
)

func (inst *Instance) Enable() error {
	inst.modbusPollingMsg("MODBUS Plugin Enable()")
	inst.enabled = true
	inst.fault = false
	inst.pluginName = name
	inst.setUUID()

	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, argspkg.Args{})
	if nets != nil {
		inst.networks = nets
	} else if err != nil {
		inst.networks = nil
	}
	if inst.config.EnablePolling {
		if !inst.pollingEnabled {
			var arg polling
			inst.pollingEnabled = true
			arg.enable = true
			inst.NetworkPollManagers = make([]*pollqueue.NetworkPollManager, 0) // This will delete any existing NetworkPollManagers (if enable is called multiple times, it will rebuild the queues).
			for _, net := range nets {                                          // Create a new Poll Manager for each network in the plugin.
				conf := inst.GetConfig().(*Config)
				if conf.PollQueueLogLevel != "ERROR" && conf.PollQueueLogLevel != "DEBUG" && conf.PollQueueLogLevel != "POLLING" {
					conf.PollQueueLogLevel = "ERROR"
				}
				pollQueueConfig := pollqueue.Config{EnablePolling: conf.EnablePolling, LogLevel: conf.PollQueueLogLevel}
				pollManager := NewPollManager(&pollQueueConfig, &inst.db, net.UUID, net.Name, inst.pluginUUID, inst.pluginName, float.NonNil(net.MaxPollRate))
				// inst.modbusDebugMsg("net")
				// inst.modbusDebugMsg("%+v\n", net)
				// inst.modbusDebugMsg("pollManager")
				// inst.modbusDebugMsg("%+v\n", pollManager)
				pollManager.StartPolling()
				inst.NetworkPollManagers = append(inst.NetworkPollManagers, pollManager)
			}
			inst.running = true
			err := inst.ModbusPolling()
			if err != nil {
				inst.fault = true
				inst.running = false
				inst.modbusErrorMsg("POLLING ERROR on routine: %v\n", err)
			}
		}
	}
	return nil
}

func (inst *Instance) Disable() error {
	inst.modbusPollingMsg("MODBUS Plugin Disable()")
	inst.enabled = false
	if inst.pollingEnabled {
		inst.pollingEnabled = false
		inst.pollingCancel()
		inst.pollingCancel = nil
		for _, pollMan := range inst.NetworkPollManagers {
			pollMan.StopPolling()
		}
		inst.NetworkPollManagers = make([]*pollqueue.NetworkPollManager, 0)
	}
	inst.running = false
	inst.fault = false
	return nil
}
