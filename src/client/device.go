package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *FlowClient) AddDevice(device *model.Device) (*model.Device, error) {
	url := fmt.Sprintf("/api/devices")
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Device{}).
		SetBody(device).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Device), nil
}

// GetFirstDevice first object
func (inst *FlowClient) GetFirstDevice(withPoints ...bool) (*model.Device, error) {
	devices, err := inst.GetDevices(withPoints...)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		return &device, err
	}
	return nil, err
}

func (inst *FlowClient) GetDevices(withPoints ...bool) ([]model.Device, error) {
	url := fmt.Sprintf("/api/devices")
	if len(withPoints) > 0 {
		if withPoints[0] == true {
			url = fmt.Sprintf("/api/devices?with_points=true")
		}
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]model.Device{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	var out []model.Device
	out = *resp.Result().(*[]model.Device)
	return out, nil
}

func (inst *FlowClient) GetDevice(uuid string, withPoints ...bool) (*model.Device, error) {
	url := fmt.Sprintf("/api/devices/%s", uuid)
	if len(withPoints) > 0 {
		if withPoints[0] == true {
			url = fmt.Sprintf("/api/devices/%s?with_points=true", uuid)
		}
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Device{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Device), nil
}

func (inst *FlowClient) EditDevice(uuid string, device *model.Device) (*model.Device, error) {
	url := fmt.Sprintf("/api/devices/%s", uuid)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Device{}).
		SetBody(device).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Device), nil
}

func (inst *FlowClient) DeleteDevice(uuid string) (bool, error) {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetPathParams(map[string]string{"uuid": uuid}).
		Delete("/api/devices/{uuid}"))
	if err != nil {
		return false, err
	}
	return true, nil
}
