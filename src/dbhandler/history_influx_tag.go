package dbhandler

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

func (h *Handler) GetHistoryInfluxTags(producerUuid string) ([]*model.HistoryInfluxTag, error) {
	return getDb().GetHistoryInfluxTags(producerUuid)
}
