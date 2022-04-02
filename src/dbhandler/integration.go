package dbhandler

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

func (h *Handler) GetEnabledIntegrationByPluginConfId(pcId string) ([]*model.Integration, error) {
	return getDb().GetEnabledIntegrationByPluginConfId(pcId)
}
