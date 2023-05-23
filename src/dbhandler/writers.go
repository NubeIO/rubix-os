package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
)

func (h *Handler) GetWriters(args api.Args) ([]*model.Writer, error) {
	q, err := getDb().GetWriters(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetWriter(uuid string) (*model.Writer, error) {
	q, err := getDb().GetWriter(uuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}
