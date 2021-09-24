package client

import (
	"fmt"
)

// ClientGetPlugins an object
func (a *FlowClient) ClientGetPlugins() (*ResponsePlugins, error) {
	resp, err := a.client.R().
		SetResult(&ResponsePlugins{}).
		Get("/plugins")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponsePlugins), nil
}

// ClientGetPlugin an object
func (a *FlowClient) ClientGetPlugin(uuid string) (*ResponseBody, error) {
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/plugins/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}
