package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
)

func (h *Handler) GetHistoriesForSync() ([]*model.History, error) {
	q, err := getDb().GetHistoriesForSync(api.Args{})
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) CreateBulkHistory(histories []*model.History) (bool, error) {
	q, err := getDb().CreateBulkHistory(histories)
	if err != nil {
		return false, err
	}
	return q, nil
}
