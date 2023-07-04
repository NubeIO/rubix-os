package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	parentArgs "github.com/NubeIO/rubix-os/args"
)

func (h *Handler) GetHosts(args parentArgs.Args) ([]*model.Host, error) {
	return getDb().GetHosts(false, args)
}
