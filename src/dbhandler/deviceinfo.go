package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetDeviceInfo() (*model.DeviceInfo, error) {
	q, err := getDb().GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	return q, nil
}
