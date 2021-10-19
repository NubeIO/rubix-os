package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/system/host"
	"github.com/NubeDev/flow-framework/src/system/networking"
	"github.com/NubeDev/flow-framework/src/system/ufw"
	"github.com/NubeDev/flow-framework/src/utilstime"
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

func (a *DeviceInfoAPI) GetSystemTime(ctx *gin.Context) {
	t := utilstime.SystemTime()
	reposeHandler(t, nil, ctx)
}

func (a *DeviceInfoAPI) GetExternalIP(ctx *gin.Context) {
	t, err := networking.ExternalIPV4()
	reposeHandler(t, err, ctx)
}

func (a *DeviceInfoAPI) GetNetworks(ctx *gin.Context) {
	_, _, all, err := networking.IpAddresses()
	reposeHandler(all, err, ctx)
}

func (a *DeviceInfoAPI) GetInterfacesNames(ctx *gin.Context) {
	t, err := networking.GetInterfacesNames()
	reposeHandler(t, err, ctx)
}

func (a *DeviceInfoAPI) GetInternetStatus(ctx *gin.Context) {
	t, err := networking.CheckInternetStatus()
	reposeHandler(t, err, ctx)
}

func (a *DeviceInfoAPI) GetOSDetails(ctx *gin.Context) {
	out := host.GetCombinationData(false)
	reposeHandler(out, nil, ctx)
}

func (a *DeviceInfoAPI) GetTZoneList(ctx *gin.Context) {
	out, err := utilstime.GetTimeZoneList()
	reposeHandler(out, err, ctx)
}

func (a *DeviceInfoAPI) FirewallStatus(ctx *gin.Context) {
	out, err := ufw.FirewallStatus()
	reposeHandler(out, err, ctx)
}
