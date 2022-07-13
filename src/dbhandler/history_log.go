package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetHistoryLogByFlowNetworkCloneUUID(fncUuid string) (*model.HistoryLog, error) {
	q, err := getDb().GetHistoryLogByFlowNetworkCloneUUID(fncUuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdateHistoryLog(log *model.HistoryLog) (*model.HistoryLog, error) {
	q, err := getDb().UpdateHistoryLog(log)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdateBulkHistoryLogs(logs []*model.HistoryLog) (bool, error) {
	return getDb().UpdateBulkHistoryLogs(logs)
}
