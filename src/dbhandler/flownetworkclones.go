package dbhandler

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetFlowNetworkClones() ([]*model.FlowNetworkClone, error) {
	q, err := getDb().GetFlowNetworkClones(api.Args{})
	if err != nil {
		return nil, err
	}
	return q, nil
}
