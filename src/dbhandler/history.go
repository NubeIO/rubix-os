package dbhandler

import (
	"github.com/NubeIO/flow-framework/model"
)

func (h *Handler) GetHistoriesForSync() ([]*model.History, error) {
	return getDb().GetHistoriesForSync()
}

func (h *Handler) CreateBulkHistory(histories []*model.History) (bool, error) {
	return getDb().CreateBulkHistory(histories)
}
