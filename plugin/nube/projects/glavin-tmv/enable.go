package main

import (
	"github.com/NubeIO/flow-framework/api"
)

func (inst *Instance) Enable() error {
	inst.tmvDebugMsg("Polling Enable()")
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
			// TODO: VERIFY POLLING WITHOUT GO ROUTINE WRAPPER
			err := inst.ModbusPolling()
			if err != nil {
				inst.tmvErrorMsg("POLLING ERROR on routine: %v\n", err)
			}
		}
	}
	return nil
}

func (inst *Instance) Disable() error {
	inst.tmvDebugMsg("Polling Disable()")
	inst.enabled = false
	if inst.pollingEnabled {
		inst.pollingEnabled = false
		inst.pollingCancel()
		inst.pollingCancel = nil
	}
	return nil
}
