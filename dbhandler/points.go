package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (l *Handler) CreatePoint(body *model.Point) (*model.Point, error) {
	q, err := getDb().CreatePoint(body)
	if err != nil {
		return nil, err
	}
	return q, nil
}
