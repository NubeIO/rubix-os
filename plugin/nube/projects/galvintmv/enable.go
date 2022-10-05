package main

import (
	"github.com/NubeIO/flow-framework/api"
)

func (inst *Instance) Enable() error {
	inst.tmvDebugMsg("TMV Plugin Enable()")
	inst.enabled = true
	inst.pluginName = name
	inst.setUUID()

	nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, api.Args{})
	if nets != nil {
		inst.networks = nets
	} else if err != nil {
		inst.networks = nil
	}
	inst.updatePointNames()
	return nil
}

func (inst *Instance) Disable() error {
	inst.tmvDebugMsg("TMV Plugin Disable()")
	inst.enabled = false
	if inst.pollingEnabled && inst.pollingCancel != nil {
		inst.pollingEnabled = false
		inst.pollingCancel()
		inst.pollingCancel = nil
	}
	return nil
}
