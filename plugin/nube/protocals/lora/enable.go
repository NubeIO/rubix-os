package main

import (
	"errors"
)

// Enable implements plugin.Plugin
func (c *Instance) Enable() error {
	//
	c.enabled = true
	c.setUUID()
	c.BusServ()
	q, err := c.db.GetNetworkByPlugin(c.pluginUUID, false, false, "serial")
	if err != nil {
		return errors.New("there is no network added please add one")
	}
	c.networkUUID = q.UUID
	err = c.SerialOpen();if err != nil {
		return errors.New("error on enable lora-plugin")

	}
	return nil
}

// Disable implements plugin.Disable
func (c *Instance) Disable() error {
	c.enabled = false
	c.SerialClose()
	return nil
}

