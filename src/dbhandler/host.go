package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetHosts() ([]*model.Host, error) {
	return getDb().GetHosts(false)
}
