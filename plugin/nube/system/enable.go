package main

import (
	"errors"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, false, false, "serial")
	if err != nil {
		return errors.New("there is no network added please add one")
	}
	i.networkUUID = q.UUID
	//err = i.SerialOpen()
	if err != nil {
		return errors.New("error on enable lora-plugin")

	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	//err := i.SerialClose()
	//if err != nil {
	//	return err
	//}
	return nil
}
