package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
)


// ClientAddSubscriber an object
func (a *FlowClient) ClientAddSubscriber(body Subscriber) (*ResponseBody, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("sub_name_%s", name)
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(body).
		Post("/api/subscriber")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	fmt.Println(resp.String())
	return resp.Result().(*ResponseBody), nil
}


// ClientGetSubscriber an object
func (a *FlowClient) ClientGetSubscriber(uuid string) (*ResponseBody, error) {
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/subscriber/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	fmt.Println(resp.Error())
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	fmt.Println(resp.String())

	return resp.Result().(*ResponseBody), nil
}


// ClientEditSubscriber edit an object
func (a *FlowClient) ClientEditSubscriber(uuid string) (*ResponseBody, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("sub_new_name_%s", name)
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(map[string]string{"name": name}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Post("/api/subscriber/{}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}

