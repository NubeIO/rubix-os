package client

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/flow-framework/utils/nstring"
)

func (inst *FlowClient) AddAutoMappings(body *interfaces.AutoMappingNetwork) *interfaces.AutoMappingNetworkError {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&interfaces.AutoMappingNetworkError{}).
		SetBody(body).
		Post("/api/auto_mappings"))
	if err != nil {
		return &interfaces.AutoMappingNetworkError{Name: body.Name, Error: nstring.New(err.Error())}
	}
	return resp.Result().(*interfaces.AutoMappingNetworkError)
}
