package client

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nresty"
)

func (inst *FlowClient) Login(body *model.LoginBody) (*model.TokenResponse, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		SetResult(&model.TokenResponse{}).
		Post("/api/users/login"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.TokenResponse), nil
}
