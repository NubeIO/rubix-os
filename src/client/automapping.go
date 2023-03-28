package client

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/nresty"
)

func (inst *FlowClient) CreateAutoMapping(body *interfaces.AutoMappingNetwork) *interfaces.AutoMappingError {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&interfaces.AutoMappingError{}).
		SetBody(body).
		Post("/api/auto_mappings"))
	if err != nil {
		return &interfaces.AutoMappingError{NetworkUUID: body.UUID, Error: err.Error(), Level: interfaces.Network}
	}
	return resp.Result().(*interfaces.AutoMappingError)
}
