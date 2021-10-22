package main

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/rubix/rubixapi"
	"github.com/labstack/gommon/log"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	var arg api.Args
	arg.WithIpConnection = true
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if q != nil {
		i.networkUUID = q.UUID
	} else {
		i.networkUUID = "NA"
	}
	if err != nil {
		log.Error("error on enable lora-plugin")
	}
	i.REST = rubixapi.New()
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}
