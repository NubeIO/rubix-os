package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/labstack/gommon/log"
)

// Enable implements plugin.Plugin
func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if q != nil {
		inst.networkUUID = q.UUID
	}
	inst.initBacStore()

	var arg polling
	inst.pollingEnabled = true
	arg.enable = true
	go func() error {
		err := inst.polling(arg)
		if err != nil {
			log.Errorf("rubix-io.enable: POLLING ERROR on routine: %v\n", err)
		}
		return nil
	}()
	if err != nil {
		log.Errorf("rubix-io.enable: POLLING ERROR: %v\n", err)
	}

	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	inst.setUUID()
	q, _ := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	inst.closeBacnetStoreNetwork(q.UUID)
	if inst.pollingEnabled {
		var arg polling
		inst.pollingEnabled = false
		arg.enable = false
		go func() {
			err := inst.polling(arg)
			if err != nil {
			}
		}()
	}
	return nil
}
