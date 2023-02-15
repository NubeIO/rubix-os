package main

import (
	"github.com/NubeIO/flow-framework/api"
)

func (inst *Instance) Enable() error {
	inst.bacnetDebugMsg("Polling Enable()")
	inst.enabled = true
	inst.running = false
	inst.fault = false
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
		inst.running = true
		inst.fault = false
		err := inst.BACnetServerPolling()
		if err != nil {
			inst.running = false
			inst.fault = true
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
	inst.running = false
	inst.fault = false
	return nil
}
