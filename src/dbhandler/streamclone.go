package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetStreamClones(args api.Args) ([]*model.StreamClone, error) {
	q, err := getDb().GetStreamClones(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}
