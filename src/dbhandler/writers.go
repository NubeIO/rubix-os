package dbhandler

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetWriters(args api.Args) ([]*model.Writer, error) {
	q, err := getDb().GetWriters(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetWriter(uuid string) (*model.Writer, error) {
	q, err := getDb().GetWriter(uuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}
