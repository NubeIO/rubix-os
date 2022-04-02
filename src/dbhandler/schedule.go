package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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

func (h *Handler) GetOneScheduleByArgs(args api.Args) (*model.Schedule, error) {
	return getDb().GetOneScheduleByArgs(args)
}

func (h *Handler) UpdateSchedule(uuid string, body *model.Schedule) (*model.Schedule, error) {
	q, err := getDb().UpdateSchedule(uuid, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}
