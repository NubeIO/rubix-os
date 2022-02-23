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

func (h *Handler) GetOneDeviceByArgs(args api.Args) (*model.Device, error) {
	return getDb().GetOneDeviceByArgs(args)
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

func (h *Handler) DeleteDevice(uuid string) (bool, error) {
	_, err := getDb().DeleteDevice(uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}
