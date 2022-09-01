package main

import (
	"github.com/NubeIO/flow-framework/api"
)

func (inst *Instance) Enable() error {
	inst.bacnetDebugMsg("Polling Enable()")
	inst.enabled = true
	inst.pluginName = name
	inst.setUUID()
	inst.BusServ()
	net, _ := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if net != nil {
		inst.networkUUID = net.UUID
	}
	inst.initBacStore()

	if !inst.pollingEnabled {
		inst.pollingEnabled = true
		err := inst.BACnetServerPolling()
		if err != nil {
			inst.bacnetErrorMsg("POLLING ERROR on routine: %v\n", err)
		}
		go inst.initPointsNames()
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
	}
	return nil
}
