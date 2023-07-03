package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
)

func (h *Handler) GetHosts(args api.Args) ([]*model.Host, error) {
	return getDb().GetHosts(false, args)
}
