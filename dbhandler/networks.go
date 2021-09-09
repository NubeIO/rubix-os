package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetNetwork(uuid string, withChildren bool, withPoints bool) (*model.Network, error) {
	q, err := getDb().GetNetwork(uuid, withChildren, withPoints)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworkByPlugin(pluginUUID string, withChildren bool, withPoints bool, transport string) (*model.Network, error) {
	q, err := getDb().GetNetworkByPlugin(pluginUUID, withChildren, withPoints, transport)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdateNetwork(uuid string, body *model.Network) (*model.Network, error) {
	q, err := getDb().UpdateNetwork(uuid, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}
