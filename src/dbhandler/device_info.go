package dbhandler

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/deviceinfo"
)

func (h *Handler) GetDeviceInfo() (*model.DeviceInfo, error) {
	q, err := deviceinfo.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	return q, nil
}
