package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
)


// ClientAddNetwork an object
func (a *FlowClient) ClientAddNetwork() (*ResponseBody, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("net_name_%s", name)
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(map[string]string{"name": name}).
		Post("/api/network")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}


// ClientGetNetwork an object
func (a *FlowClient) ClientGetNetwork(uuid string) (*ResponseBody, error) {
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/network/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}
