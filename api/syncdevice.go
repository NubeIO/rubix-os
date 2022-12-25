package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type SyncDeviceDatabase interface {
	SyncDevice(fn *model.SyncDevice) (*model.Device, error)
}

type SyncDeviceAPI struct {
	DB SyncDeviceDatabase
}

func (a *SyncDeviceAPI) SyncDevice(ctx *gin.Context) {
	body, _ := getBodySyncDevice(ctx)
	q, err := a.DB.SyncDevice(body)
	ResponseHandler(q, err, ctx)
}
