package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetSchedules() ([]*model.Schedule, error) {
	q, err := getDb().GetSchedules()
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetSchedule(uuid string) (*model.Schedule, error) {
	q, err := getDb().GetSchedule(uuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdateSchedule(uuid string, body *model.Schedule) (*model.Schedule, error) {
	q, err := getDb().UpdateSchedule(uuid, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}
