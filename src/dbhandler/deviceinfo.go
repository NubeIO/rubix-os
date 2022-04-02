package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetDeviceInfo() (*model.DeviceInfo, error) {
	q, err := getDb().GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	return q, nil
}
