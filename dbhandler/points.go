package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetPoint(uuid string, withChildren bool) (*model.Point, error) {
	q, err := getDb().GetPoint(uuid, withChildren)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) CreatePoint(body *model.Point) (*model.Point, error) {
	q, err := getDb().CreatePoint(body)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdatePoint(uuid string, body *model.Point) (*model.Point, error) {
	q, err := getDb().UpdatePoint(uuid, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}
