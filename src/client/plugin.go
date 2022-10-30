package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// ClientGetPlugins an object
func (inst *FlowClient) ClientGetPlugins() (*ResponsePlugins, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&ResponsePlugins{}).
		Get("/plugins"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponsePlugins), nil
}

// ClientGetPlugin an object
func (inst *FlowClient) ClientGetPlugin(uuid string) (*ResponseBody, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/plugins/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*ResponseBody), nil
}

// CreateNetworkPlugin an object
func (inst *FlowClient) CreateNetworkPlugin(body *model.Network, pluginName string) (*model.Network, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/networks", pluginName)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Network{}).
		SetBody(body).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Network), nil
}

// DeleteNetworkPlugin delete an object
func (inst *FlowClient) DeleteNetworkPlugin(body *model.Network, pluginName string) (ok bool, err error) {
	url := fmt.Sprintf("/api/plugins/api/%s/networks", pluginName)
	_, err = nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Delete(url))
	if err != nil {
		return false, err
	}
	return true, err
}

// DeleteDevicePlugin delete an object
func (inst *FlowClient) DeleteDevicePlugin(body *model.Device, pluginName string) (ok bool, err error) {
	url := fmt.Sprintf("/api/plugins/api/%s/devices", pluginName)
	_, err = nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Delete(url))
	if err != nil {
		return false, err
	}
	return true, err
}

// DeletePointPlugin delete an object
func (inst *FlowClient) DeletePointPlugin(body *model.Point, pluginName string) (ok bool, err error) {
	url := fmt.Sprintf("/api/plugins/api/%s/points", pluginName)
	_, err = nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Delete(url))
	if err != nil {
		return false, err
	}
	return true, err
}

// CreateDevicePlugin an object
func (inst *FlowClient) CreateDevicePlugin(body *model.Device, pluginName string) (*model.Device, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/devices", pluginName)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Device{}).
		SetBody(body).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Device), nil
}

// CreatePointPlugin an object
func (inst *FlowClient) CreatePointPlugin(body *model.Point, pluginName string) (*model.Point, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/points", pluginName)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Point{}).
		SetBody(body).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}

// UpdateNetworkPlugin update an object
func (inst *FlowClient) UpdateNetworkPlugin(body *model.Network, pluginName string) (*model.Network, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/networks", pluginName)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Network{}).
		SetBody(body).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Network), nil
}

// UpdateDevicePlugin update an object
func (inst *FlowClient) UpdateDevicePlugin(body *model.Device, pluginName string) (*model.Device, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/devices", pluginName)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Device{}).
		SetBody(body).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Device), nil
}

// UpdatePointPlugin update an object
func (inst *FlowClient) UpdatePointPlugin(body *model.Point, pluginName string) (*model.Point, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/points", pluginName)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Point{}).
		SetBody(body).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}

// WritePointPlugin update an object
func (inst *FlowClient) WritePointPlugin(pointUUID string, body *model.PointWriter, pluginName string) (*model.Point, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/points/write/{uuid}", pluginName)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Point{}).
		SetBody(body).
		SetPathParams(map[string]string{"uuid": pointUUID}).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}
