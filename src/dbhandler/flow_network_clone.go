package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetFlowNetworkClones(args api.Args) ([]*model.FlowNetworkClone, error) {
	q, err := getDb().GetFlowNetworkClones(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) RefreshFlowNetworkClonesConnections() (*bool, error) {
	return getDb().RefreshFlowNetworkClonesConnections()
}
