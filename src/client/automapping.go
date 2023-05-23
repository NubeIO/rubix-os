package client

import (
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *FlowClient) CreateAutoMapping(body *interfaces.AutoMapping) interfaces.AutoMappingResponse {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&interfaces.AutoMappingResponse{}).
		SetBody(body).
		Post("/api/auto_mappings"))
	if err != nil {
		networkUUID := "" // pick first valid network
		for _, network := range body.Networks {
			if network.CreateNetwork {
				networkUUID = network.UUID
				break
			}
		}
		return interfaces.AutoMappingResponse{
			HasError:    true,
			NetworkUUID: networkUUID,
			Error:       err.Error(),
			Level:       interfaces.Network,
		}
	}
	return *resp.Result().(*interfaces.AutoMappingResponse)
}
