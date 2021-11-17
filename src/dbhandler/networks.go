package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
)

func (h *Handler) CreateNetwork(body *model.Network) (*model.Network, error) {
	q, err := getDb().CreateNetwork(body)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworkByPlugin(pluginUUID string, args api.Args) (*model.Network, error) {
	q, err := getDb().GetNetworkByPlugin(pluginUUID, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworksByPlugin(pluginUUID string, args api.Args) ([]*model.Network, error) {
	q, err := getDb().GetNetworksByPlugin(pluginUUID, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworks(args api.Args) ([]*model.Network, error) {
	q, err := getDb().GetNetworks(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworksByName(name string, args api.Args) ([]*model.Network, error) {
	q, err := getDb().GetNetworksByName(name, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworkByField(field string, value string, withDevices bool) (*model.Network, error) {
	q, err := getDb().GetNetworkByField(field, value, withDevices)
	if err != nil {
		return nil, err
	}
	return q, nil
}
