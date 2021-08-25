package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
)


// ClientAddGateway an object
func (a *FlowClient) ClientAddGateway(isRemote bool) (*ResponseBody, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("gte_name_%s", name)
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(Stream{"name", isRemote}).
		Post("/api/stream")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}


// ClientGetGateway an object
func (a *FlowClient) ClientGetGateway(uuid string) (*ResponseBody, error) {
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/stream/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}


// ClientEditGateway edit an object
func (a *FlowClient) ClientEditGateway(uuid string) (*ResponseBody, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("dev_new_name_%s", name)
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(map[string]string{"name": name}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Post("/api/stream/{}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}

