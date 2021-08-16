package api

import (
	"github.com/NubeDev/plug-framework/helpers"
	"github.com/NubeDev/plug-framework/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// The NetworkDatabase interface for encapsulating database access.
type NetworkDatabase interface {
	GetNetwork(uuid string) (*model.Network, error)
	GetNetworks() ([]*model.Network, error)
	CreateNetwork(network *model.Network) error
	UpdateNetwork(uuid string, body *model.Network) (*model.Network, error)
	DeleteNetwork(uuid string) (bool, error)

}
type NetworksAPI struct {
	DB       NetworkDatabase
}

func (a *NetworksAPI) GetNetworks(ctx *gin.Context) {
	apps, err := a.DB.GetNetworks()
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	ctx.JSON(http.StatusOK, apps)

}

func (a *NetworksAPI) GetNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	apps, err := a.DB.GetNetwork(uuid)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	ctx.JSON(http.StatusOK, apps)

}


func (a *NetworksAPI) UpdateNetwork(ctx *gin.Context) {
	body, _ := getBODY(ctx)
	uuid := resolveID(ctx)
	apps, err := a.DB.UpdateNetwork(uuid, body)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	ctx.JSON(http.StatusOK, apps)

}

func (a *NetworksAPI) CreateNetwork(ctx *gin.Context) {
	app := model.Network{}
	app.Uuid, _ = helpers.MakeUUID()
	if err := ctx.Bind(&app); err == nil {
		if success := successOrAbort(ctx, http.StatusInternalServerError, a.DB.CreateNetwork(&app)); !success {
			return
		}
		ctx.JSON(http.StatusOK, app)
	}
}


func (a *NetworksAPI) DeleteNetwork(ctx *gin.Context) {
	uuid := resolveID(ctx)
	apps, err := a.DB.DeleteNetwork(uuid)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	ctx.JSON(http.StatusOK, apps)

}

