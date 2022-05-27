package client

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (a *FlowClient) DeviceInfo() (*model.DeviceInfo, error) {
	resp, err := FormatRestyResponse(a.client.R().
		SetResult(&model.DeviceInfo{}).
		Get("/api/system/device_info"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.DeviceInfo), nil
}
