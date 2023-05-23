package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *FlowClient) AddNetwork(body *model.Network) (*model.Network, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Network{}).
		SetBody(body).
		Post("/api/networks"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Network), nil
}

func (inst *FlowClient) EditNetwork(uuid string, body *model.Network) (*model.Network, error) {
	url := fmt.Sprintf("/api/networks/%s", uuid)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Network{}).
		SetBody(body).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Network), nil
}

func (inst *FlowClient) DeleteNetwork(uuid string) (bool, error) {
	url := fmt.Sprintf("/api/networks/%s", uuid)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		Delete(url))
	if err != nil {
		return false, err
	}
	if resp.IsSuccess() {
		return true, nil
	}
	return false, nil
}

func (inst *FlowClient) GetNetworkByPluginName(pluginName string, withPoints ...bool) (*model.Network, error) {
	url := fmt.Sprintf("/api/networks/plugin/%s", pluginName)
	if len(withPoints) > 0 {
		url = fmt.Sprintf("/api/networks/plugin/%s?with_points=true", pluginName)
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Network{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Network), nil
}

func (inst *FlowClient) GetNetworksWithPoints() ([]model.Network, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]model.Network{}).
		Get("/api/networks?with_points=true"))
	if err != nil {
		return nil, err
	}
	var out []model.Network
	out = *resp.Result().(*[]model.Network)
	return out, nil
}

func (inst *FlowClient) GetNetworkWithPoints(uuid string) (*model.Network, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Network{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/networks/{uuid}?with_points=true"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Network), nil
}

func (inst *FlowClient) GetFirstNetwork(withDevices ...bool) (*model.Network, error) {
	nets, err := inst.GetNetworks(withDevices...)
	if err != nil {
		return nil, err
	}
	for _, net := range nets {
		return &net, err
	}
	return nil, err
}

func (inst *FlowClient) GetNetworks(withDevices ...bool) ([]model.Network, error) {
	url := fmt.Sprintf("/api/networks")
	if len(withDevices) > 0 {
		if withDevices[0] == true {
			url = fmt.Sprintf("/api/networks?with_devices=true")
		}
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]model.Network{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	var out []model.Network
	out = *resp.Result().(*[]model.Network)
	return out, nil
}

// GetNetwork an object
func (inst *FlowClient) GetNetwork(uuid string) (*model.Network, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Network{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/networks/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Network), nil
}

func (inst *FlowClient) GetNetworkV2(uuid string) (*model.Network, error, error) {
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.client.R().
		SetResult(&model.Network{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/networks/{uuid}"))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*model.Network), nil, nil
}

func (inst *FlowClient) GetNetworkByName(networkName string) (*model.Network, error, error) {
	url := fmt.Sprintf("/api/networks/name/%s", networkName)
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.client.R().
		SetResult(&model.Network{}).
		Get(url))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*model.Network), nil, nil
}
