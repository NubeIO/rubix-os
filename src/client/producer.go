package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// AddProducer an object
func (a *FlowClient) AddProducer(body model.Producer) (*model.Producer, error) {
	name, _ := nuuid.MakeUUID()
	name = fmt.Sprintf("sub_name_%s", name)
	resp, err := a.client.R().
		SetResult(&model.Producer{}).
		SetBody(body).
		Post("/api/producers")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.Producer), nil
}

func (a *FlowClient) GetProducers(streamUUID *string) (*[]model.Producer, error) {
	req := a.client.R().
		SetResult(&[]model.Producer{})
	if streamUUID != nil {
		req.SetQueryParam("stream_uuid", *streamUUID)
	}
	resp, err := req.
		Get("/api/producers")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("GetProducers: %s", err)
		} else {
			return nil, fmt.Errorf("GetProducers: %s", resp)
		}
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetProducers: %s", resp)
	}
	return resp.Result().(*[]model.Producer), nil
}

func (a *FlowClient) GetProducer(uuid string) (*model.Producer, error) {
	resp, err := a.client.R().
		SetResult(&model.Producer{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/producers/{uuid}")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("GetProducer: %s", err)
		} else {
			return nil, fmt.Errorf("GetProducer: %s", resp)
		}
	}
	if resp.IsError() {
		return nil, fmt.Errorf("GetProducer: %s", resp)
	}
	return resp.Result().(*model.Producer), nil
}

// EditProducer edit an object
func (a *FlowClient) EditProducer(uuid string, body model.Producer) (*model.Producer, error) {
	resp, err := a.client.R().
		SetResult(&model.Producer{}).
		SetBody(body).
		SetPathParams(map[string]string{"uuid": uuid}).
		Patch("/api/producers/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	fmt.Println(resp.String())
	return resp.Result().(*model.Producer), nil
}
