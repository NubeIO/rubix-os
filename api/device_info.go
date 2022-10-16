package api

import (
	"github.com/NubeIO/flow-framework/utils/deviceinfo"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type DeviceInfoDatabase interface {
	GetDeviceInfo() (*model.DeviceInfo, error)
}

type DeviceInfoAPI struct {
}

func (inst *DeviceInfoAPI) GetDeviceInfo(ctx *gin.Context) {
	q, err := deviceinfo.GetDeviceInfo()
	ResponseHandler(q, err, ctx)
}
