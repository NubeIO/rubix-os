package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type DeviceInfoDatabase interface {
	GetDeviceInfo() (*model.DeviceInfo, error)
}

type DeviceInfoAPI struct {
	DB DeviceInfoDatabase
}

func (a *DeviceInfoAPI) GetDeviceInfo(ctx *gin.Context) {
	q, err := a.DB.GetDeviceInfo()
	reposeHandler(q, err, ctx)
}
