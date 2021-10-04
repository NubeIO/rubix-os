package main

import (
	"github.com/NubeDev/flow-framework/api"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	var arg api.Args
	arg.SerialConnection = true
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if q != nil {
		i.networkUUID = q.UUID
	} else {
		i.networkUUID = "NA"
	}
	if err == nil {
		i.networkUUID = q.UUID
		err = i.SerialOpen()
		if err != nil {
			//return errors.New("error on enable lora-plugin")
		}
	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	err := i.SerialClose()
	if err != nil {
		return err
	}
	return nil
}
