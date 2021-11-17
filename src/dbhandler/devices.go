package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
)

func (h *Handler) GetDevice(uuid string, args api.Args) (*model.Device, error) {
	q, err := getDb().GetDevice(uuid, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetDeviceByField(field string, value string, withPoints bool) (*model.Device, error) {
	q, err := getDb().GetDeviceByField(field, value, withPoints)
	if err != nil {
		return nil, err
	}
	return q, nil
}

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
