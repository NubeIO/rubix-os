package client

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/nresty"
)

func (inst *FlowClient) AddAutoMapping(body *interfaces.AutoMapping) (bool, error) {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Post("/api/auto_mappings"))
	if err != nil {
		return false, err
	}
	return true, nil
}
