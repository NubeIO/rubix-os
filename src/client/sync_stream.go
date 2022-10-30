package client

import (
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *FlowClient) SyncStream(body *model.SyncStream) (*model.StreamClone, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.StreamClone{}).
		SetBody(body).
		Post("/api/sync/stream"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.StreamClone), nil
}
