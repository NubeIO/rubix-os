package dbhandler

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

func (h *Handler) GetDeviceMetaTags() ([]*model.DeviceMetaTag, error) {
	return getDb().GetDeviceMetaTags()
}
