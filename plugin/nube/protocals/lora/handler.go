package main

import "github.com/NubeDev/flow-framework/model"

func (i *Instance) network(withChildren bool, withPoints bool, transport string) (*model.Network, error) {
	q, err := i.db.GetNetworkByPlugin(i.pluginUUID, withChildren, withPoints, transport)
	if err != nil {
		return nil, err
	}
	i.networkUUID = q.UUID
	return q, err
}
