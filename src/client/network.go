package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/utils/nuuid"
)

// ClientAddNetwork an object
func (a *FlowClient) ClientAddNetwork(pluginUUID string) (*ResponseBody, error) {
	name, _ := nuuid.MakeUUID()
	name = fmt.Sprintf("net_name_%s", name)
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(map[string]string{"name": name, "plugin_conf_id": pluginUUID}).
		Post("/api/networks")
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
		Get("/api/networks/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}
