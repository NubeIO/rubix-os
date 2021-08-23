package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The DeviceDatabase interface for encapsulating database access.
type DeviceDatabase interface {
	GetDevice(uuid string, withPoints bool) (*model.Device, error)
	GetDevices(withPoints bool) ([]*model.Device, error)
	CreateDevice(body *model.Device) (*model.Device, error)
	UpdateDevice(uuid string, body *model.Device) (*model.Device, error)
	DeleteDevice(uuid string) (bool, error)
	DropDevices() (bool, error)

}
type DeviceAPI struct {
	DB DeviceDatabase
}

func (a *DeviceAPI) GetDevices(ctx *gin.Context) {
	q, err := a.DB.GetDevices(false)
	reposeHandler(q, err, ctx)

}

func (a *DeviceAPI) GetDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetDevice(uuid, false)
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
