package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/interfaces"
)

func (h *Handler) CreateNetwork(body *model.Network) (*model.Network, error) {
	q, err := getDb().CreateNetwork(body)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// UpdateNetworkErrors will only update the error properties of the network, all other properties will not be updated.
func (h *Handler) UpdateNetworkErrors(uuid string, body *model.Network) error {
	err := getDb().UpdateNetworkErrors(uuid, body)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) UpdateNetwork(uuid string, body *model.Network) (*model.Network, error) {
	return getDb().UpdateNetwork(uuid, body)
}

func (h *Handler) GetNetwork(uuid string, args argspkg.Args) (*model.Network, error) {
	q, err := getDb().GetNetwork(uuid, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworkByPlugin(pluginUUID string, args argspkg.Args) (*model.Network, error) {
	q, err := getDb().GetNetworkByPlugin(pluginUUID, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworksByPluginName(pluginUUID string, args argspkg.Args) ([]*model.Network, error) {
	q, err := getDb().GetNetworksByPluginName(pluginUUID, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworksByPlugin(pluginUUID string, args argspkg.Args) ([]*model.Network, error) {
	q, err := getDb().GetNetworksByPlugin(pluginUUID, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworks(args argspkg.Args) ([]*model.Network, error) {
	q, err := getDb().GetNetworks(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworkByDeviceUUID(uuid string, args argspkg.Args) (*model.Network, error) {
	q, err := getDb().GetNetworkByDeviceUUID(uuid, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworkByName(name string, args argspkg.Args) (*model.Network, error) {
	q, err := getDb().GetNetworkByName(name, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) DeleteNetwork(uuid string) (bool, error) {
	_, err := getDb().DeleteNetwork(uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *Handler) SetErrorsForAllDevicesOnNetwork(networkUUID string, message string, messageLevel string, messageCode string, doPoints bool) error {
	err := getDb().SetErrorsForAllDevicesOnNetwork(networkUUID, message, messageLevel, messageCode, doPoints)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) ClearErrorsForAllDevicesOnNetwork(networkUUID string, doPoints bool) error {
	err := getDb().ClearErrorsForAllDevicesOnNetwork(networkUUID, doPoints)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetNetworksTagsForPostgresSync() ([]*interfaces.NetworkTagForPostgresSync, error) {
	return getDb().GetNetworksTagsForPostgresSync()
}

func (h *Handler) GetOneNetworkByArgs(args argspkg.Args) (*model.Network, error) {
	return getDb().GetOneNetworkByArgs(args)
}
