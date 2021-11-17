package main

import (
	"github.com/NubeIO/flow-framework/api"
	log "github.com/sirupsen/logrus"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	var arg api.Args
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if q != nil {
		i.networkUUID = q.UUID
	} else {
		i.networkUUID = "NA"
	}
	if err != nil {
		log.Error("system-plugin: error on enable system-plugin")
	}
	i.schedule()
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}
