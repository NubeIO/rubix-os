package client

import (
	"fmt"
)

type Ping struct {
	Health   string `json:"health"`
	Database string `json:"database"`
}

func (a *FlowClient) Ping() (*Ping, error) {
	resp, err := a.client.R().
		SetResult(&Ping{}).
		Get("/api/system/ping")
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("ping: %s", err)
		} else {
			return nil, fmt.Errorf("ping: %s", resp)
		}
	}
	return resp.Result().(*Ping), nil
}
