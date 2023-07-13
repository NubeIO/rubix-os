package client

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *FlowClient) GetScheduleV2(uuid string) (*model.Schedule, error, error) {
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.client.R().
		SetResult(&model.Schedule{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/schedules/{uuid}"))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*model.Schedule), nil, nil
}

func (inst *FlowClient) GetScheduleByNameV2(name string) (*model.Schedule, error, error) {
	resp, connectionErr, requestErr := nresty.FormatRestyV2Response(inst.client.R().
		SetResult(&model.Schedule{}).
		SetQueryParam("name", name).
		Get("/api/schedules/one/args"))
	if connectionErr != nil || requestErr != nil {
		return nil, connectionErr, requestErr
	}
	return resp.Result().(*model.Schedule), nil, nil
}
