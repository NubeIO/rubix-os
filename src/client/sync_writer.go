package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (inst *FlowClient) SyncWriter(body *model.SyncWriter) (*model.WriterClone, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.WriterClone{}).
		SetBody(body).
		Post("/api/sync/writer"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.WriterClone), nil
}

func (inst *FlowClient) SyncCOV(writerUUID string, body *model.SyncCOV) error {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Post(fmt.Sprintf("/api/sync/cov/%s", writerUUID)))
	return err
}

func (inst *FlowClient) SyncWriterWriteAction(sourceUUID string, body *model.SyncWriterAction) error {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		Post(fmt.Sprintf("/api/sync/writer/write/%s", sourceUUID)))
	return err
}

func (inst *FlowClient) SyncWriterReadAction(sourceUUID string) error {
	_, err := nresty.FormatRestyResponse(inst.client.R().Get(fmt.Sprintf("/api/sync/writer/read/%s", sourceUUID)))
	return err
}
