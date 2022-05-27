package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// ClientAddGateway an object
func (a *FlowClient) ClientAddGateway(body *model.Stream) (*ResponseBody, error) {
	name, _ := nuuid.MakeUUID()
	name = fmt.Sprintf("gte_name_%s", name)
	resp, err := CheckError(a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(body).
		Post("/api/streams"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}

// ClientGetGateway an object
func (a *FlowClient) ClientGetGateway(uuid string) (*ResponseBody, error) {
	resp, err := CheckError(a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/streams/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}

// ClientEditGateway edit an object
func (a *FlowClient) ClientEditGateway(uuid string) (*ResponseBody, error) {
	name, _ := nuuid.MakeUUID()
	name = fmt.Sprintf("dev_new_name_%s", name)
	resp, err := CheckError(a.client.R().
		SetResult(&ResponseBody{}).
		SetBody(map[string]string{"name": name}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Post("/api/streams/{}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}
