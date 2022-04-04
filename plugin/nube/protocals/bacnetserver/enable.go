package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/api/bacnetserver/v1/bsrest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/api/common/v1/iorest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/rest/v1/rest"
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

	restService := &rest.Service{}
	restService.Port = 1717
	options := &rest.Options{}
	restService.Options = options
	restService = rest.New(restService)

	commonClient := &iorest.NubeRest{}
	commonClient.UseRubixProxy = false
	commonClient.RubixUsername = "admin"
	commonClient.RubixPassword = "N00BWires"
	commonClient = iorest.New(commonClient, restService)
	bacnetClient = bsrest.New(&bsrest.BacnetClient{IoRest: commonClient})

	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	return nil
}
