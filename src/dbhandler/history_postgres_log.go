package dbhandler

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

func (h *Handler) GetHistoryPostgresLogLastSyncHistoryId() (int, error) {
	q, err := getDb().GetHistoryPostgresLogLastSyncHistoryId()
	if err != nil {
		return 0, err
	}
	return q, nil
}

func (h *Handler) UpdateHistoryPostgresLog(log *model.HistoryPostgresLog) (*model.HistoryPostgresLog, error) {
	q, err := getDb().UpdateHistoryPostgresLog(log)
	if err != nil {
		return nil, err
	}
	return q, nil
}
