package client

import (
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *FlowClient) SyncFlowNetwork(body *model.FlowNetwork) (*model.FlowNetworkClone, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.FlowNetworkClone{}).
		SetBody(body).
		Post("/api/sync/flow_network"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.FlowNetworkClone), nil
}
