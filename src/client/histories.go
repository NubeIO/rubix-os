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
		Get("/api/histories/producers/latest/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("%s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*model.ProducerHistory), nil
}
