package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetIntegrations() ([]*model.Integration, error) {
	q, err := getDb().GetIntegrations()
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetIntegration(uuid string) (*model.Integration, error) {
	q, err := getDb().GetIntegration(uuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetIntegrationByName(name string) (*model.Integration, error) {
	q, err := getDb().GetIntegrationByName(name)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdateIntegration(uuid string, body *model.Integration) (*model.Integration, error) {
	q, err := getDb().UpdateIntegration(uuid, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}
