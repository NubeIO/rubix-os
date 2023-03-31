package client

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/nresty"
)

func (inst *FlowClient) CreateAutoMapping(body *interfaces.AutoMapping) interfaces.AutoMappingResponse {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&interfaces.AutoMappingResponse{}).
		SetBody(body).
		Post("/api/auto_mappings"))
	if err != nil {
		return interfaces.AutoMappingResponse{HasError: true, Error: err.Error()}
	}
	return *resp.Result().(*interfaces.AutoMappingResponse)
}
