package dbhandler

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetProducerHistories() ([]*model.ProducerHistory, error) {
	q, err := getDb().GetProducerHistories(api.Args{})
	if err != nil {
		return nil, err
	}
	return q, nil
}
