package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nresty"
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

func (inst *FlowClient) GetDeviceV2(uuid string) (*model.Device, error, error) {
	url := fmt.Sprintf("/api/devices/%s", uuid)
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.client.R().
		SetResult(&model.Device{}).
		Get(url))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*model.Device), nil, nil
}

func (inst *FlowClient) GetDeviceByName(networkName, deviceName string) (*model.Device, error, error) {
	url := fmt.Sprintf("/api/devices/name/%s/%s", networkName, deviceName)
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.client.R().
		SetResult(&model.Device{}).
		Get(url))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*model.Device), nil, nil
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
