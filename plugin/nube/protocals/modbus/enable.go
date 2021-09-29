package main

import (
	"errors"
	"github.com/NubeDev/flow-framework/api"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	var arg api.Args
	arg.IpConnection = true
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, arg)
	if err != nil {
		return errors.New("there is no network added please add one")
	}
	i.networkUUID = q.UUID
	if err != nil {
		return errors.New("error on enable lora-plugin")
	}
	if i.config.EnablePolling {
		if !i.pollingEnabled {
			var arg polling
			i.pollingEnabled = true
			arg.enable = true
			go func() {
				err := i.PollingTCP(arg)
				if err != nil {

				}
			}()
			if err != nil {
				return errors.New("error on starting polling")
			}
		}
	}
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false

	if i.pollingEnabled {
		var arg polling
		i.pollingEnabled = false
		arg.enable = false
		go func() {
			err := i.PollingTCP(arg)
			if err != nil {

			}
		}()
		if err != nil {
			return errors.New("error on starting polling")
		}
	}

	return nil
}
