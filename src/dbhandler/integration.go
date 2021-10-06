package dbhandler

import "github.com/NubeDev/flow-framework/model"

func (h *Handler) GetEnabledIntegrationByPluginConfId(pcId string) ([]*model.Integration, error) {
	q, err := getDb().GetEnabledIntegrationByPluginConfId(pcId)
	if err != nil {
		return nil, err
	}
	return q, nil
}
