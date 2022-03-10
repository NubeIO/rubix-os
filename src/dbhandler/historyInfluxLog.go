package dbhandler

import (
	"github.com/NubeIO/flow-framework/model"
)

func (h *Handler) UpdateHistoryInfluxLogLastSyncId(lastSyncId int) (*model.HistoryInfluxLog, error) {
	q, err := getDb().UpdateHistoryInfluxLogLastSyncId(lastSyncId)
	if err != nil {
		return nil, err
	}
	return q, nil
}
