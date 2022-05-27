package client

import "github.com/NubeIO/flow-framework/nresty"

type Ping struct {
	Health   string `json:"health"`
	Database string `json:"database"`
}

func (a *FlowClient) Ping() (*Ping, error) {
	resp, err := nresty.FormatRestyResponse(a.client.R().
		SetResult(&Ping{}).
		Get("/api/system/ping"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*Ping), nil
}
