package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/labstack/gommon/log"
)

// Enable implements plugin.Plugin
func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.BusServ()
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if q != nil {
		inst.networkUUID = q.UUID
	} else {
		inst.networkUUID = "NA"
	}
	if err == nil {
		inst.networkUUID = q.UUID
		err = inst.SerialOpen()
		if err != nil {
			log.Error("error on enable lora-plugin")
		}
	}
	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	err := inst.SerialClose()
	if err != nil {
		return err
	}
	return nil
}
