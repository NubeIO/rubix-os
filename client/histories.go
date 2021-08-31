package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
)



// ClientGetHistory an object
func (a *FlowClient) ClientGetHistory(uuid string) (*model.ProducerHistory, error) {
	resp, err := a.client.R().
		SetResult(&model.ProducerHistory{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/histories/by/producer/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.ProducerHistory), nil
}


