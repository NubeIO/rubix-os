package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/api/bacnetserver/v1/bsrest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/api/common/v1/iorest"
	"github.com/labstack/gommon/log"
)

var bacnetClient *bsrest.BacnetClient

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

	commonClient := new(iorest.NubeRest)
	commonClient.UseRubixProxy = true
	commonClient.RubixUsername = "admin"
	commonClient.RubixPassword = "N00BWires"
	commonClient = iorest.New(commonClient)
	//bacnetClient.Url = "0.0.0.0"
	//bacnetClient.Port = 1717
	//bacnetClient.IoRest = commonClient
	bacnetClient = bsrest.New(&bsrest.BacnetClient{Url: "0.0.0.0", Port: 1717, IoRest: commonClient})

	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	return nil
}
