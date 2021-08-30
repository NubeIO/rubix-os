package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


// ClientAddPoint an object
func (a *FlowClient) ClientAddPoint(deviceUUID string) (*model.Point, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("pnt_name_%s", name)
	resp, err := a.client.R().
		SetResult(&model.Point{}).
		SetBody(map[string]string{"name": name, "device_uuid": deviceUUID}).
		Post("/api/point")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.Point), nil
}


// ClientGetPoint an object
func (a *FlowClient) ClientGetPoint(uuid string) (*model.Point, error) {
	resp, err := a.client.R().
		SetResult(&model.Point{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/point/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.Point), nil
}


// ClientEditPoint an object
func (a *FlowClient) ClientEditPoint(uuid string, body model.Point) (*model.Point, error) {
	resp, err := a.client.R().
		SetBody(body).
		SetResult(&model.Point{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Patch("/api/point/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.Point), nil
}


