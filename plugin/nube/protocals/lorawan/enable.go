package main

import (
	"errors"
	rest "github.com/NubeDev/flow-framework/plugin/nube/protocals/lorawan/restclient"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, false, false, "ip")
	if err != nil {
		return errors.New("there is no network added please add one")
	}
	i.networkUUID = q.UUID
	if err != nil {
		return errors.New("error on enable lora-plugin")
	}
	i.CLI = rest.NewChirp(chirpName, chirpPass, ip, port)
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}
