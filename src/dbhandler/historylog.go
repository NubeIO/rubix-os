package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) UpdateHistoryLogLastSyncId(lastSyncId int) (*model.HistoryLog, error) {
	q, err := getDb().UpdateHistoryLogLastSyncId(lastSyncId)
	if err != nil {
		return nil, err
	}
	return q, nil
}
