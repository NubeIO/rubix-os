package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type PointDatabase interface {
	GetPoints(args Args) ([]*model.Point, error)
	GetPoint(uuid string, args Args) (*model.Point, error)
	CreatePoint(body *model.Point, fromPlugin bool) (*model.Point, error)
	UpdatePoint(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error)
	PointWrite(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error)
	GetOnePointByArgs(args Args) (*model.Point, error)
	DeletePoint(uuid string) (bool, error)
	DropPoints() (bool, error)
}
type PointAPI struct {
	DB PointDatabase
}

func (a *PointAPI) GetPoints(ctx *gin.Context) {
	args := buildPointArgs(ctx)
	q, err := a.DB.GetPoints(args)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) GetPoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildPointArgs(ctx)
	q, err := a.DB.GetPoint(uuid, args)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) UpdatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdatePoint(uuid, body, false)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) PointWrite(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.PointWrite(uuid, body, false)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) GetOnePointByArgs(ctx *gin.Context) {
	args := buildPointArgs(ctx)
	q, err := a.DB.GetOnePointByArgs(args)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) CreatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	q, err := a.DB.CreatePoint(body, false)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) DeletePoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeletePoint(uuid)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) DropPoints(ctx *gin.Context) {
	q, err := a.DB.DropPoints()
	responseHandler(q, err, ctx)
}
