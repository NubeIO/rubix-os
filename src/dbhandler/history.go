package dbhandler

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetHistoriesForSync() ([]*model.History, error) {
	q, err := getDb().GetHistoriesForSync(api.Args{})
	if err != nil {
		return nil, err
	}
	return q, nil
}
