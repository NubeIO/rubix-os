package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/rubixregistry"
	"github.com/gin-gonic/gin"
)

type DeviceInfoDatabase interface {
	GetDeviceInfo() (*model.DeviceInfo, error)
}

type DeviceInfoAPI struct {
	RubixRegistry *rubixregistry.RubixRegistry
}

func (a *DeviceInfoAPI) GetDeviceInfo(c *gin.Context) {
	deviceInfo, err := a.RubixRegistry.GetDeviceInfo()
	ResponseHandler(deviceInfo, err, c)
}
