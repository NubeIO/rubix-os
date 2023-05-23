package client

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/nresty"
	"strconv"
)

// GetWriterClone an object
func (inst *FlowClient) GetWriterClone(uuid string) (*model.WriterClone, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.WriterClone{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/producers/writer_clones/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.WriterClone), nil
}

// EditWriterClone edit an object
func (inst *FlowClient) EditWriterClone(uuid string, body model.WriterClone, updateProducer bool) (*model.WriterClone, error) {
	param := strconv.FormatBool(updateProducer)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.WriterClone{}).
		SetBody(body).
		SetPathParams(map[string]string{"uuid": uuid}).
		SetQueryParam("update_producer", param).
		Patch("/api/producers/writer_clones/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.WriterClone), nil
}

// CreateWriterClone edit an object
func (inst *FlowClient) CreateWriterClone(body model.WriterClone) (*model.WriterClone, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.WriterClone{}).
		SetBody(body).
		Post("/api/producers/writer_clones"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.WriterClone), nil
}
