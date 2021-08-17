package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// The PointDatabase interface for encapsulating database access.
type PointDatabase interface {
	GetPoint(uuid string, withChildren bool) (*model.Point, error)
	GetPoints(withChildren bool) ([]*model.Point, error)
	CreatePoint(points *model.Point, body *model.Point) error
	UpdatePoint(uuid string, body *model.Point) (*model.Point, error)
	DeletePoint(uuid string) (bool, error)
}
type PointAPI struct {
	DB PointDatabase
}

func (a *PointAPI) GetPoints(ctx *gin.Context) {
	q, err := a.DB.GetPoints(false)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}

func (a *PointAPI) GetPoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetPoint(uuid, false)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}

func (a *PointAPI) UpdatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdatePoint(uuid, body)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}

func (a *PointAPI) CreatePoint(ctx *gin.Context) {
	app := model.Point{}
	body, _ := getBODYPoint(ctx)
	if success := successOrAbort(ctx, http.StatusInternalServerError, a.DB.CreatePoint(&app, body)); !success {
		return
	}
	ctx.JSON(http.StatusOK, app)

}

func (a *PointAPI) DeletePoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeletePoint(uuid)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}
