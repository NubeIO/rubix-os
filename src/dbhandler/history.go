package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetHistoriesForSync(lastSyncId int) ([]*model.History, error) {
	return getDb().GetHistoriesForSync(lastSyncId)
}

func (h *Handler) CreateBulkHistory(histories []*model.History) (bool, error) {
	return getDb().CreateBulkHistory(histories)
}
