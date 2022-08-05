package main

import (
	"github.com/NubeIO/flow-framework/api"
)

func (inst *Instance) Enable() error {
	inst.edge28DebugMsg("Edge28 Enable()")
	inst.enabled = true
	inst.pluginName = name
	inst.setUUID()
	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
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

			err := inst.Edge28Polling()
			if err != nil {
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
	return nil
}
