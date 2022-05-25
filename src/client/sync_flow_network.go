package client

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (a *FlowClient) SyncFlowNetwork(body *model.FlowNetwork) (*model.FlowNetworkClone, error) {
	resp, err := a.client.R().
		SetResult(&model.FlowNetworkClone{}).
		SetBody(body).
		Post("/api/sync/flow_network")
	fr := failedResponse(err, resp)
	if fr != nil {
		return nil, err
	}
	return resp.Result().(*model.FlowNetworkClone), nil
}
