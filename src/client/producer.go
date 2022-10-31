package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// AddProducer an object
func (inst *FlowClient) AddProducer(body *model.Producer) (*model.Producer, error) {
	name, _ := nuuid.MakeUUID()
	name = fmt.Sprintf("sub_name_%s", name)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Producer{}).
		SetBody(body).
		Post("/api/producers"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Producer), nil
}

func (inst *FlowClient) GetProducers() ([]model.Producer, error) {
	url := fmt.Sprintf("/api/producers")
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]model.Producer{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	var out []model.Producer
	out = *resp.Result().(*[]model.Producer)
	return out, nil
}

func (inst *FlowClient) GetProducer(uuid string) (*model.Producer, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Producer{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/producers/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Producer), nil
}

// EditProducer edit an object
func (inst *FlowClient) EditProducer(uuid string, body *model.Producer) (*model.Producer, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Producer{}).
		SetBody(body).
		SetPathParams(map[string]string{"uuid": uuid}).
		Patch("/api/producers/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Producer), nil
}

func (inst *FlowClient) DeleteProducer(uuid string) (bool, error) {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetPathParams(map[string]string{"uuid": uuid}).
		Delete("/api/producers/{uuid}"))
	if err != nil {
		return false, err
	}
	return true, nil
}
