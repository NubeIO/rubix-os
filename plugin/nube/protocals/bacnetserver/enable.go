package main

import (
	"github.com/NubeIO/flow-framework/api"
	nube_api_bacnetserver "github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/api/nube_api/bacnetserver"
	"github.com/labstack/gommon/log"
)

var bacnetClient = nube_api_bacnetserver.New("", 0)

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
	if err != nil {
		log.Error("error on enable bacnetserver-plugin")
	}
	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	return nil
}
