package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) CloneEdge(host *model.Host) error {
	return getDb().CloneEdge(host)
}
