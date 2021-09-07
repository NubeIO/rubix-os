package main

// Enable implements plugin.Plugin
func (c *Instance) Enable() error {
	//go SerialOpenAndRead()
	c.enabled = true
	c.setUUID()
	return nil
}

// Disable implements plugin.Disable
func (c *Instance) Disable() error {
	c.enabled = false
	return nil
}

