package dbhandler

import (
	"github.com/NubeIO/flow-framework/model"
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
