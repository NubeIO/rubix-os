package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *FlowClient) AddPoint(body *model.Point) (*model.Point, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Point{}).
		SetBody(body).
		Post("/api/points"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}

func (inst *FlowClient) GetPoints() ([]model.Point, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]model.Point{}).
		Get("/api/points"))
	if err != nil {
		return nil, err
	}
	var out []model.Point
	out = *resp.Result().(*[]model.Point)
	return out, nil
}

func (inst *FlowClient) GetPoint(uuid string) (*model.Point, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Point{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/points/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}

func (inst *FlowClient) GetPointsBulkUUIDs() (*[]string, error, error) {
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.client.R().
		SetResult(&[]string{}).
		Get("/api/points/bulk/uuids"))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*[]string), nil, nil
}

func (inst *FlowClient) GetPointByName(networkName, deviceName, pointName string) (*model.Point, error, error) {
	url := fmt.Sprintf("/api/points/name/%s/%s/%s", networkName, deviceName, pointName)
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.client.R().
		SetResult(&model.Point{}).
		Get(url))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*model.Point), nil, nil
}

func (inst *FlowClient) DeletePoint(uuid string) (bool, error) {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetPathParams(map[string]string{"uuid": uuid}).
		Delete("/api/points/{uuid}"))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (inst *FlowClient) EditPoint(uuid string, body *model.Point) (*model.Point, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		SetResult(&model.Point{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Patch("/api/points/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}

func (inst *FlowClient) GetPointWithParent(uuid string) (*interfaces.PointWithParent, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&interfaces.PointWithParent{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/points/with_parent/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*interfaces.PointWithParent), nil
}

func (inst *FlowClient) WritePoint(uuid string, body *model.PointWriter) (*model.Point, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		SetResult(&model.Point{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Patch("/api/points/write/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}
