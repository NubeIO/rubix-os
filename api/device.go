package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type DeviceDatabase interface {
	GetDevices(args Args) ([]*model.Device, error)
	GetDevice(uuid string, args Args) (*model.Device, error)
	CreateDevice(body *model.Device) (*model.Device, error)
	UpdateDevice(uuid string, body *model.Device) (*model.Device, error)
	DeleteDevice(uuid string) (bool, error)
	DropDevices() (bool, error)
	GetDeviceByField(field string, value string, withPoints bool) (*model.Device, error)
	UpdateDeviceByField(field string, value string, body *model.Device) (*model.Device, error)
}
type DeviceAPI struct {
	DB DeviceDatabase
}

func (a *DeviceAPI) GetDevices(ctx *gin.Context) {
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetDevices(args)
	reposeHandler(q, err, ctx)
}

func (a *DeviceAPI) GetDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetDevice(uuid, args)
	reposeHandler(q, err, ctx)
}

func (a *DeviceAPI) GetDeviceByField(ctx *gin.Context) {
	field, value := withFieldsArgs(ctx)
	q, err := a.DB.GetDeviceByField(field, value, false)
	reposeHandler(q, err, ctx)
}

func (a *DeviceAPI) UpdateDeviceByField(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	field, value := withFieldsArgs(ctx)
	q, err := a.DB.UpdateDeviceByField(field, value, body)
	reposeHandler(q, err, ctx)
}

func (a *DeviceAPI) UpdateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateDevice(uuid, body)
	reposeHandler(q, err, ctx)
}

func (a *DeviceAPI) CreateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	q, err := a.DB.CreateDevice(body)
	reposeHandler(q, err, ctx)
}

func (a *DeviceAPI) DeleteDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteDevice(uuid)
	reposeHandler(q, err, ctx)
}

func (a *DeviceAPI) DropDevices(ctx *gin.Context) {
	q, err := a.DB.DropDevices()
	reposeHandler(q, err, ctx)
}
