package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
)

func (h *Handler) GetStreamClones(args api.Args) ([]*model.StreamClone, error) {
	q, err := getDb().GetStreamClones(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}
