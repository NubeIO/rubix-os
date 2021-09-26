package dbhandler

import "github.com/NubeDev/flow-framework/model"

func (h *Handler) GetProducerHistories() ([]*model.ProducerHistory, error) {
	q, err := getDb().GetProducerHistories()
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetProducerHistory(uuid string) (*model.ProducerHistory, error) {
	q, err := getDb().GetProducerHistory(uuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}
