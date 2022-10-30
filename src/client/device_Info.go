package client

import (
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *FlowClient) DeviceInfo() (*model.DeviceInfo, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.DeviceInfo{}).
		Get("/api/system/device_info"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.DeviceInfo), nil
}
