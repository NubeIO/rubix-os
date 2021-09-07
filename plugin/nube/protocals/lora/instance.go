package main

func (c *Instance) setUUID()  {
	q, err := c.db.GetPluginByPath(name)
	if err != nil {
		return
	}
	c.pluginUUID = q.UUID
}

