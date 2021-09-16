package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) CreateDevice(body *model.Device) (*model.Device, error) {
	q, err := getDb().CreateDevice(body)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdateDevice(uuid string, body *model.Device) (*model.Device, error) {
	q, err := getDb().UpdateDevice(uuid, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}
