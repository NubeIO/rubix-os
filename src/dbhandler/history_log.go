package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetHistoryLogByHostUUID(hostUUID string) (*model.HistoryLog, error) {
	return getDb().GetHistoryLogByHostUUID(hostUUID)
}

func (h *Handler) UpdateHistoryLog(log *model.HistoryLog) (*model.HistoryLog, error) {
	return getDb().UpdateHistoryLog(log)
}

func (h *Handler) UpdateBulkHistoryLogs(logs []*model.HistoryLog) (bool, error) {
	return getDb().UpdateBulkHistoryLogs(logs)
}
