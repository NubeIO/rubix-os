package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// ClientAddGateway an object
func (inst *FlowClient) ClientAddGateway(body *model.Stream) (*ResponseBody, error) {
	name, _ := nuuid.MakeUUID()
	name = fmt.Sprintf("gte_name_%s", name)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&ResponseBody{}).
		SetBody(body).
		Post("/api/streams"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}

// ClientGetGateway an object
func (inst *FlowClient) ClientGetGateway(uuid string) (*ResponseBody, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/streams/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}

// ClientEditGateway edit an object
func (inst *FlowClient) ClientEditGateway(uuid string) (*ResponseBody, error) {
	name, _ := nuuid.MakeUUID()
	name = fmt.Sprintf("dev_new_name_%s", name)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&ResponseBody{}).
		SetBody(map[string]string{"name": name}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Post("/api/streams/{}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}
