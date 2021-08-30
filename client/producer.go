package client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
)


// ClientAddProducer an object
func (a *FlowClient) ClientAddProducer(body model.Producer) (*model.Producer, error) {
	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("sub_name_%s", name)
	resp, err := a.client.R().
		SetResult(&model.Producer{}).
		SetBody(body).
		Post("/api/producer")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	fmt.Println(resp.String())
	return resp.Result().(*model.Producer), nil
}


// ClientGetProducer an object
func (a *FlowClient) ClientGetProducer(uuid string) (*model.Producer, error) {
	resp, err := a.client.R().
		SetResult(&model.Producer{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/producer/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	fmt.Println(resp.Error())
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	fmt.Println(resp.String())

	return resp.Result().(*model.Producer), nil
}


// ClientEditProducer edit an object
func (a *FlowClient) ClientEditProducer(uuid string, body model.Producer) (*model.Producer, error) {
	resp, err := a.client.R().
		SetResult(&model.Producer{}).
		SetBody(body).
		SetPathParams(map[string]string{"uuid": uuid}).
		Patch("/api/producer/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	fmt.Println(resp.String())
	return resp.Result().(*model.Producer), nil
}

