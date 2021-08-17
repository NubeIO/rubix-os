package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// The DeviceDatabase interface for encapsulating database access.
type DeviceDatabase interface {
	GetDevice(uuid string, withPoints bool) (*model.Device, error)
	GetDevices(withPoints bool) ([]*model.Device, error)
	CreateDevice(device *model.Device, body *model.Device) error
	UpdateDevice(uuid string, body *model.Device) (*model.Device, error)
	DeleteDevice(uuid string) (bool, error)
}
type DeviceAPI struct {
	DB DeviceDatabase
}

func (a *DeviceAPI) GetDevices(ctx *gin.Context) {
	q, err := a.DB.GetDevices(false)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}

func (a *DeviceAPI) GetDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetDevice(uuid, false)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}

func (a *DeviceAPI) UpdateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateDevice(uuid, body)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}

func (a *DeviceAPI) CreateDevice(ctx *gin.Context) {
	app := model.Device{}
	body, _ := getBODYDevice(ctx)
	if success := successOrAbort(ctx, http.StatusInternalServerError, a.DB.CreateDevice(&app, body)); !success {
		return
	}
	ctx.JSON(http.StatusOK, app)

}

func (a *DeviceAPI) DeleteDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteDevice(uuid)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}
