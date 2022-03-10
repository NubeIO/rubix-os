package dbhandler

import "github.com/NubeIO/flow-framework/model"

func (h *Handler) GetHistoryInfluxTags(producerUuid string) ([]*model.HistoryInfluxTag, error) {
	q, err := getDb().GetHistoryInfluxTags(producerUuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}
