package main

import (
	"github.com/NubeIO/flow-framework/api"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	var arg api.Args
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, arg)
	if q != nil {
		inst.networkUUID = q.UUID
	} else {
		inst.networkUUID = "NA"
	}
	if err != nil {
		log.Error("system-plugin: error on enable system-plugin")
	}
	inst.schedule()
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	return nil
}
