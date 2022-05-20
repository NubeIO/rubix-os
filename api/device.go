package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type DeviceDatabase interface {
	GetDevices(args Args) ([]*model.Device, error)
	GetDevice(uuid string, args Args) (*model.Device, error)
	GetOneDeviceByArgs(args Args) (*model.Device, error)
	CreateDevice(body *model.Device) (*model.Device, error)
	UpdateDevice(uuid string, body *model.Device, fromPlugin bool) (*model.Device, error)
	DeleteDevice(uuid string) (bool, error)

	CreateDevicePlugin(body *model.Device) (*model.Device, error)
	UpdateDevicePlugin(uuid string, body *model.Device) (*model.Device, error)
	DeleteDevicePlugin(uuid string) (bool, error)
}
type DeviceAPI struct {
	DB DeviceDatabase
}

func (a *DeviceAPI) GetDevices(ctx *gin.Context) {
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetDevices(args)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) GetDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetDevice(uuid, args)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) GetOneDeviceByArgs(ctx *gin.Context) {
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetOneDeviceByArgs(args)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) UpdateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateDevicePlugin(uuid, body)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) CreateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	q, err := a.DB.CreateDevicePlugin(body)
	responseHandler(q, err, ctx)
}

func (a *DeviceAPI) DeleteDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteDevicePlugin(uuid)
	responseHandler(q, err, ctx)
}
