package dbhandler

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

func (h *Handler) GetLatestProducerHistoryByProducerName(name string) (*model.ProducerHistory, error) {
	q, err := getDb().GetLatestProducerHistoryByProducerName(name)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetProducerHistoriesByProducerName(name string) ([]*model.ProducerHistory, int64, error) {
	q, count, err := getDb().GetProducerHistoriesByProducerName(name)
	if err != nil {
		return nil, 0, err
	}
	return q, count, nil
}
