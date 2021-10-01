package main

import (
	"errors"
	"github.com/NubeDev/flow-framework/api"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	var arg api.Args
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if err != nil {
		return errors.New("there is no network added please add one")
	}
	i.networkUUID = q.UUID
	if err != nil {
		return errors.New("error on enable lora-plugin")

	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}
