package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The PointDatabase interface for encapsulating database access.
type PointDatabase interface {
	GetPoint(uuid string, withChildren bool) (*model.Point, error)
	GetPoints(withChildren bool) ([]*model.Point, error)
	CreatePoint(body *model.Point) (*model.Point, error)
	UpdatePoint(uuid string, body *model.Point) (*model.Point, error)
	DeletePoint(uuid string) (bool, error)
	DropPoints() (bool, error)
}
type PointAPI struct {
	DB PointDatabase
}

func (a *PointAPI) GetPoints(ctx *gin.Context) {
	withChildren, _ := withChildrenArgs(ctx)
	q, err := a.DB.GetPoints(withChildren)
	reposeHandler(q, err, ctx)

}

func (a *PointAPI) GetPoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	withChildren, _ := withChildrenArgs(ctx)
	q, err := a.DB.GetPoint(uuid, withChildren)
	reposeHandler(q, err, ctx)

}

func (a *PointAPI) UpdatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdatePoint(uuid, body)
	reposeHandler(q, err, ctx)

}

func (a *PointAPI) CreatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	q, err := a.DB.CreatePoint(body)
	reposeHandler(q, err, ctx)

}

func (a *PointAPI) DeletePoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeletePoint(uuid)
	reposeHandler(q, err, ctx)
}

func (a *PointAPI) DropPoints(ctx *gin.Context) {
	q, err := a.DB.DropPoints()
	reposeHandler(q, err, ctx)
}

