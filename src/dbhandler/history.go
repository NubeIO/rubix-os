package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"time"
)

func (h *Handler) GetHistoriesForPostgresSync(lastSyncId int) ([]*model.History, error) {
	return getDb().GetHistoriesForPostgresSync(lastSyncId)
}

func (h *Handler) GetHistoriesByHostUUID(hostUUID string, startTime, endTime time.Time) ([]*model.History, error) {
	return getDb().GetHistoriesByHostUUID(hostUUID, startTime, endTime)
}

func (h *Handler) CreateBulkHistory(histories []*model.History) (bool, error) {
	return getDb().CreateBulkHistory(histories)
}
