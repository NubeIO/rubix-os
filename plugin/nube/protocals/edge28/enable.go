package main

import (
	argspkg "github.com/NubeIO/rubix-os/args"
)

func (inst *Instance) Enable() error {
	inst.edge28DebugMsg("Edge28 Enable()")
	inst.enabled = true
	inst.running = false
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
			inst.running = true
			err := inst.Edge28Polling()
			if err != nil {
				inst.fault = true
				inst.running = false
				inst.edge28ErrorMsg("POLLING ERROR on routine: %v\n", err)
			}
		}
	}
	return nil
}

func (inst *Instance) Disable() error {
	inst.edge28DebugMsg("EDGE28 Disable()")
	inst.enabled = false
	if inst.pollingEnabled {
		var arg polling
		inst.pollingEnabled = false
		arg.enable = false
		inst.pollingCancel()
		inst.pollingCancel = nil
	}
	inst.running = false
	inst.fault = false
	return nil
}
