package client

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *FlowClient) DeviceInfo() (*model.DeviceInfo, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.DeviceInfo{}).
		Get("/api/system/device"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.DeviceInfo), nil
}
