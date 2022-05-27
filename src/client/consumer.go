package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/utils/nuuid"
)

// ClientAddConsumer an object
func (a *FlowClient) ClientAddConsumer(body Consumer) (*ResponseBody, error) {
	name, _ := nuuid.MakeUUID()
	name = fmt.Sprintf("sub_name_%s", name)
	resp, err := CheckError(a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(body).
		Post("/api/consumers"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}

// ClientGetConsumer an object
func (a *FlowClient) ClientGetConsumer(uuid string) (*ResponseBody, error) {
	resp, err := CheckError(a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/consumers/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}

// ClientEditConsumer edit an object
func (a *FlowClient) ClientEditConsumer(uuid string) (*ResponseBody, error) {
	name, _ := nuuid.MakeUUID()
	name = fmt.Sprintf("sub_new_name_%s", name)
	resp, err := CheckError(a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(map[string]string{"name": name}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Post("/api/consumers/{}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}
