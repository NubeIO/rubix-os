package main

import (
	"github.com/NubeIO/rubix-os/api"
)

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if q != nil {
		inst.networkUUID = q.UUID
	}
	inst.interruptChan = make(chan struct{}, 1)
	if err == nil {
		go inst.run()
		inst.running = true
		inst.fault = false
	} else {
		inst.running = false
		inst.fault = true
	}
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	select {
	case inst.interruptChan <- struct{}{}:
	default:
	}
	inst.running = false
	inst.fault = false
	return nil
}
