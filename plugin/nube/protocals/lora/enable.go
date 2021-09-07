package main

// Enable implements plugin.Plugin
func (c *Instance) Enable() error {
	//go SerialOpenAndRead()
	c.enabled = true

	c.setUUID()
	c.BusServ()
	q, err := c.db.GetNetworkByPlugin(c.pluginUUID, false, false, "serial")
	if err != nil {
		return err
	}
	c.networkUUID = q.UUID
	return nil
}

// Disable implements plugin.Disable
func (c *Instance) Disable() error {
	c.enabled = false
	return nil
}

