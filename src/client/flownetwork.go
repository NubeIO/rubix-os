package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *FlowClient) AddFlowNetwork(body *model.FlowNetwork) (*model.FlowNetwork, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.FlowNetwork{}).
		SetBody(body).
		Post("/api/flow_networks"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.FlowNetwork), nil
}

func (inst *FlowClient) EditFlowNetwork(uuid string, body *model.FlowNetwork) (*model.FlowNetwork, error) {
	url := fmt.Sprintf("/api/flow_networks/%s", uuid)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.FlowNetwork{}).
		SetBody(body).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.FlowNetwork), nil
}

func (inst *FlowClient) GetFlowNetworks(withStreams ...bool) ([]model.FlowNetwork, error) {
	url := fmt.Sprintf("/api/flow_networks")
	if len(withStreams) > 0 {
		if withStreams[0] == true {
			url = fmt.Sprintf("/api/flow_networks?with_streams=true")
		}
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]model.FlowNetwork{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	var out []model.FlowNetwork
	out = *resp.Result().(*[]model.FlowNetwork)
	return out, nil
}

func (inst *FlowClient) GetFlowNetwork(uuid string, withStreams ...bool) (*model.FlowNetwork, error) {
	url := fmt.Sprintf("/api/flow_networks/%s", uuid)
	if len(withStreams) > 0 {
		if withStreams[0] == true {
			url = fmt.Sprintf("/api/flow_networks/%s?with_streams=true", uuid)
		}
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.FlowNetwork{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.FlowNetwork), nil
}

func (inst *FlowClient) GetFlowNetworksWithChild() ([]model.FlowNetwork, error) {
	url := fmt.Sprintf("/api/flow_networks?with_streams=true&with_producers=true&with_consumers=true&with_writers=true&with_tags=true")
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]model.FlowNetwork{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	var out []model.FlowNetwork
	out = *resp.Result().(*[]model.FlowNetwork)
	return out, nil
}

func (inst *FlowClient) GetFlowNetworkWithChild(uuid string) (*model.FlowNetwork, error) {
	url := fmt.Sprintf("/api/flow_networks/%s?with_streams=true&with_producers=true&with_consumers=true&with_writers=true&with_tags=true", uuid)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.FlowNetwork{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.FlowNetwork), nil
}

func (inst *FlowClient) DeleteFlowNetwork(uuid string) (bool, error) {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetPathParams(map[string]string{"uuid": uuid}).
		Delete("/api/flow_networks/{uuid}"))
	if err != nil {
		return false, err
	}
	return true, nil
}
