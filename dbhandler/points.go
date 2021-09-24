package dbhandler

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetPoints(args api.Args) ([]*model.Point, error) {
	q, err := getDb().GetPoints(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPoint(uuid string, withChildren bool) (*model.Point, error) {
	q, err := getDb().GetPoint(uuid, api.Args{})
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) CreatePoint(body *model.Point) (*model.Point, error) {
	q, err := getDb().CreatePoint(body, "")
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdatePoint(uuid string, body *model.Point, writeValue, fromPlugin bool) (*model.Point, error) {
	q, err := getDb().UpdatePoint(uuid, body, writeValue, fromPlugin)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetPointByField(field string, value string, withChildren bool) (*model.Point, error) {
	q, err := getDb().GetPointByField(field, value, withChildren)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdatePointByFieldAndType(field string, value string, body *model.Point) (*model.Point, error) {
	q, err := getDb().UpdatePointByFieldAndType(field, value, body, false)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) PointAndQuery(value1, value2 string) (*model.Point, error) {
	q, err := getDb().PointAndQuery(value1, value2)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) DeletePoint(uuid string) (bool, error) {
	_, err := getDb().DeletePoint(uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}
