package client

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/nresty"
)

func (inst *FlowClient) CreateAutoMapping(body *interfaces.AutoMappingNetwork) interfaces.AutoMappingResponse {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&interfaces.AutoMappingResponse{}).
		SetBody(body).
		Post("/api/auto_mappings"))
	if err != nil {
		return interfaces.AutoMappingResponse{NetworkUUID: body.UUID, Error: err.Error(), Level: interfaces.Network}
	}
	return *resp.Result().(*interfaces.AutoMappingResponse)
}
