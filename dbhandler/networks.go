package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetNetworkByPlugin(pluginUUID string, withChildren bool, withPoints bool, transport string) (*model.Network, error) {
	q, err := getDb().GetNetworkByPlugin(pluginUUID, withChildren, withPoints, transport)
	if err != nil {
		return nil, err
	}
	return q, nil
}
