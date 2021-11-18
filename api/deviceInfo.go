package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/system/host"
	"github.com/NubeIO/flow-framework/src/system/networking"
	"github.com/NubeIO/flow-framework/src/system/ufw"
	"github.com/NubeIO/flow-framework/src/utilstime"
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
	responseHandler(q, err, ctx)
}

func (a *DeviceInfoAPI) GetSystemTime(ctx *gin.Context) {
	t := utilstime.SystemTime()
	responseHandler(t, nil, ctx)
}

func (a *DeviceInfoAPI) GetExternalIP(ctx *gin.Context) {
	t, err := networking.ExternalIPV4()
	responseHandler(t, err, ctx)
}

func (a *DeviceInfoAPI) GetNetworks(ctx *gin.Context) {
	_, _, all, err := networking.IpAddresses()
	responseHandler(all, err, ctx)
}

func (a *DeviceInfoAPI) GetInterfacesNames(ctx *gin.Context) {
	t, err := networking.GetInterfacesNames()
	responseHandler(t, err, ctx)
}

func (a *DeviceInfoAPI) GetInternetStatus(ctx *gin.Context) {
	t, err := networking.CheckInternetStatus()
	responseHandler(t, err, ctx)
}

func (a *DeviceInfoAPI) GetOSDetails(ctx *gin.Context) {
	out := host.GetCombinationData(false)
	responseHandler(out, nil, ctx)
}

func (a *DeviceInfoAPI) GetTZoneList(ctx *gin.Context) {
	out, err := utilstime.GetTimeZoneList()
	responseHandler(out, err, ctx)
}

func (a *DeviceInfoAPI) FirewallStatus(ctx *gin.Context) {
	out, err := ufw.FirewallStatus()
	responseHandler(out, err, ctx)
}
