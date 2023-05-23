package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
)

func (h *Handler) GetConsumers(args api.Args) ([]*model.Consumer, error) {
	q, err := getDb().GetConsumers(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}
