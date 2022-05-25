package client

type Ping struct {
	Health   string `json:"health"`
	Database string `json:"database"`
}

func (a *FlowClient) Ping() (*Ping, error) {
	resp, err := a.client.R().
		SetResult(&Ping{}).
		Get("/api/system/ping")
	fr := failedResponse(err, resp)
	if fr != nil {
		return nil, err
	}
	return resp.Result().(*Ping), nil
}
