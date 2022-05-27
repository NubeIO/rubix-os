package client

import (
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (a *FlowClient) Login(body *model.LoginBody) (*model.Token, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetBody(body).
		SetResult(&model.Token{}).
		Post("/api/users/login"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Token), nil
}
