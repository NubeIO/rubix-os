package main

import "github.com/NubeDev/flow-framework/model"

func (c *Instance) network(withChildren bool, withPoints bool, transport string) (*model.Network, error)  {
	q, err := c.db.GetNetworkByPlugin(c.pluginUUID, withChildren, withPoints, transport);if err != nil {
		return nil, err
	}
	c.networkUUID = q.UUID
	return q, err
}