package client

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (a *FlowClient) Login(body *model.LoginBody) (*model.Token, error) {
	resp, err := a.client.R().
		SetBody(body).
		SetResult(&model.Token{}).
		Post("/api/users/login")
	err = CheckError(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Token), nil
}
