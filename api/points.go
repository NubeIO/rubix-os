package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/gin-gonic/gin"
)

// The PointDatabase interface for encapsulating database access.
type PointDatabase interface {
	GetPoints(args Args) ([]*model.Point, error)
	GetPoint(uuid string, args Args) (*model.Point, error)
	GetPointsByNetworkPluginName(networkUUID string) (*utils.Array, error)
	GetPointsByNetworkUUID(networkUUID string) (*utils.Array, error)
	CreatePoint(body *model.Point, addToParent string) (*model.Point, error)
	UpdatePoint(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error)
	PointWrite(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error)
	GetPointByName(networkName, deviceName, pointName string) (*model.Point, error)
	GetPointByField(field string, value string) (*model.Point, error)
	PointWriteByName(networkName, deviceName, pointName string, body *model.Point, fromPlugin bool) (*model.Point, error)
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

func (a *PointAPI) GetPointsByNetworkPluginName(ctx *gin.Context) {
	name := resolveName(ctx)
	q, err := a.DB.GetPointsByNetworkPluginName(name)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointsByNetworkUUID(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetPointsByNetworkUUID(uuid)
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

func (a *PointAPI) GetPointByName(ctx *gin.Context) {
	networkName, deviceName, pointName := networkDevicePointNames(ctx)
	q, err := a.DB.GetPointByName(networkName, deviceName, pointName)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointByField(ctx *gin.Context) {
	field, value := withFieldsArgs(ctx)
	q, err := a.DB.GetPointByField(field, value)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) PointWriteByName(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	networkName, deviceName, pointName := networkDevicePointNames(ctx)
	q, err := a.DB.PointWriteByName(networkName, deviceName, pointName, body, false)
	responseHandler(q, err, ctx)
}

func (a *PointAPI) CreatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	addToParent := parentArgs(ctx) //flowUUID
	q, err := a.DB.CreatePoint(body, addToParent)
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
