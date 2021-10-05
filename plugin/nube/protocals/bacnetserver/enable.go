package main

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/labstack/gommon/log"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	var arg api.Args
	arg.IpConnection = true
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if q != nil {
		i.networkUUID = q.UUID
	} else {
		i.networkUUID = "NA"
	}
	if err != nil {
		log.Error("error on enable bacnetserver-plugin")
	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}
