package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
)


// ClientAddSubscription an object
func (a *FlowClient) ClientAddSubscription(body Subscription) (*ResponseBody, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("sub_name_%s", name)
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(body).
		Post("/api/subscription")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	fmt.Println(resp.String())
	return resp.Result().(*ResponseBody), nil
}


// ClientGetSubscription an object
func (a *FlowClient) ClientGetSubscription(uuid string) (*ResponseBody, error) {
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/subscription/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}


// ClientEditSubscription edit an object
func (a *FlowClient) ClientEditSubscription(uuid string) (*ResponseBody, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("sub_new_name_%s", name)
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(map[string]string{"name": name}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Post("/api/subscription/{}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}

