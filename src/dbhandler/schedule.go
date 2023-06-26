package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetSchedules() ([]*model.Schedule, error) {
	q, err := getDb().GetSchedules()
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdateScheduleAllProps(uuid string, body *model.Schedule) (*model.Schedule, error) {
	q, err := getDb().UpdateScheduleAllProps(uuid, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}
