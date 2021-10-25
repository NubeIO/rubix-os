package dbhandler

import "github.com/NubeDev/flow-framework/model"

func (h *Handler) GetWriters() ([]*model.Writer, error) {
	q, err := getDb().GetWriters()
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

func (h *Handler) GetWritersByThingClass(writerThingClass string) ([]*model.Writer, error) {
	q, err := getDb().GetWritersByThingClass(writerThingClass)
	if err != nil {
		return nil, err
	}
	return q, nil
}
