package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
)

func (a *FlowClient) DeviceInfo() (*model.DeviceInfo, error) {
	resp, err := a.client.R().
		SetResult(&model.DeviceInfo{}).
		Get("/api/system/device_info")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("DeviceInfo: %s", err)
		} else {
			return nil, fmt.Errorf("DeviceInfo: %s", resp)
		}
	}
	return resp.Result().(*model.DeviceInfo), nil
}
