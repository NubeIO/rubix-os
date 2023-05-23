package client

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nresty"
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
