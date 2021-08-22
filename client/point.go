package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
)


// ClientAddPoint an object
func (a *FlowClient) ClientAddPoint(deviceUUID string) (*ResponseBody, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("net_%s", name)
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(map[string]string{"name": name, "device_uuid": deviceUUID}).
		Post("/api/point")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}


// ClientGetPoint an object
func (a *FlowClient) ClientGetPoint(uuid string) (*ResponseBody, error) {
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/point/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}
