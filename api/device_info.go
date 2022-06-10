package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type DeviceInfoDatabase interface {
	GetDeviceInfo() (*model.DeviceInfo, error)
}

type DeviceInfoAPI struct {
	DB DeviceInfoDatabase
}

func (inst *DeviceInfoAPI) GetDeviceInfo(ctx *gin.Context) {
	q, err := inst.DB.GetDeviceInfo()
	ResponseHandler(q, err, ctx)
}
