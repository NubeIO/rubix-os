package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
)

func (a *FlowClient) SyncFlowNetwork(body *model.FlowNetwork) (*model.FlowNetworkClone, error) {
	resp, err := a.client.R().
		SetResult(&model.FlowNetworkClone{}).
		SetBody(body).
		Post("/api/sync/flow_network")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("SyncFlowNetwork: %s", err)
		} else {
			return nil, fmt.Errorf("SyncFlowNetwork: %s", resp)
		}
	}
	if resp.IsError() {
		return nil, fmt.Errorf("SyncFlowNetwork: %s", resp)
	}
	return resp.Result().(*model.FlowNetworkClone), nil
}
