package api

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// The DeviceDatabase interface for encapsulating database access.
type DeviceDatabase interface {
	GetDevice(uuid string) (*model.Device, error)
	GetDevices() ([]*model.Device, error)
	CreateDevice(device *model.Device, body *model.Device) error
	UpdateDevice(uuid string, body *model.Device) (*model.Device, error)
	DeleteDevice(uuid string) (bool, error)
}
type DeviceAPI struct {
	DB DeviceDatabase
}

func (a *DeviceAPI) GetDevices(ctx *gin.Context) {
	apps, err := a.DB.GetDevices()
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	ctx.JSON(http.StatusOK, apps)

}

func (a *DeviceAPI) GetDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	apps, err := a.DB.GetDevice(uuid)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	ctx.JSON(http.StatusOK, apps)

}

func (a *DeviceAPI) UpdateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	uuid := resolveID(ctx)
	apps, err := a.DB.UpdateDevice(uuid, body)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	ctx.JSON(http.StatusOK, apps)

}

func (a *DeviceAPI) CreateDevice(ctx *gin.Context) {
	app := model.Device{}
	body, _ := getBODYDevice(ctx)
	fmt.Println(body, 66666)
	if success := successOrAbort(ctx, http.StatusInternalServerError, a.DB.CreateDevice(&app, body)); !success {
		return
	}
	ctx.JSON(http.StatusOK, app)

}

func (a *DeviceAPI) DeleteDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	apps, err := a.DB.DeleteDevice(uuid)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	ctx.JSON(http.StatusOK, apps)

}
