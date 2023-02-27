package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type PointDatabase interface {
	GetPoints(args Args) ([]*model.Point, error)
	GetPointsBulk(bulkPoints []*model.Point) ([]*model.Point, error)
	GetPoint(uuid string, args Args) (*model.Point, error)
	CreatePoint(body *model.Point, fromPlugin bool) (*model.Point, error)
	UpdatePoint(uuid string, body *model.Point, fromPlugin bool, afterRealDeviceUpdate bool) (*model.Point, error)
	PointWrite(uuid string, body *model.PointWriter, afterRealDeviceUpdate bool, currentWriterUUID *string, forceWrite bool) (*model.Point, bool, bool, bool, error)
	GetOnePointByArgs(args Args) (*model.Point, error)
	DeletePoint(uuid string) (bool, error)
	GetPointByName(networkName, deviceName, pointName string, args Args) (*model.Point, error)
	PointWriteByName(networkName, deviceName, pointName string, body *model.PointWriter) (*model.Point, error)
	DeleteOnePointByArgs(args Args) (bool, error)

	CreatePointPlugin(body *model.Point) (*model.Point, error)
	UpdatePointPlugin(uuid string, body *model.Point) (*model.Point, error)
	WritePointPlugin(uuid string, body *model.PointWriter) (*model.Point, error)
	DeletePointPlugin(uuid string) (bool, error)

	CreatePointMetaTags(pointUUID string, pointMetaTags []*model.PointMetaTag) ([]*model.PointMetaTag, error)
}
type PointAPI struct {
	DB PointDatabase
}

func (a *PointAPI) GetPoints(ctx *gin.Context) {
	args := buildPointArgs(ctx)
	q, err := a.DB.GetPoints(args)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointsBulk(ctx *gin.Context) {
	body, _ := getBODYBulkPoints(ctx)
	q, err := a.DB.GetPointsBulk(body)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) GetPoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildPointArgs(ctx)
	q, err := a.DB.GetPoint(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) UpdatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdatePointPlugin(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) PointWrite(ctx *gin.Context) {
	body, _ := getBODYPointWriter(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.WritePointPlugin(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) GetOnePointByArgs(ctx *gin.Context) {
	args := buildPointArgs(ctx)
	q, err := a.DB.GetOnePointByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) CreatePoint(ctx *gin.Context) {
	body, _ := getBODYPoint(ctx)
	q, err := a.DB.CreatePointPlugin(body)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) DeletePoint(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeletePointPlugin(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointByNameArgs(ctx *gin.Context) {
	networkName, deviceName, pointName := networkDevicePointNames(ctx)
	q, err := a.DB.GetPointByName(networkName, deviceName, pointName, Args{})
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointByName(ctx *gin.Context) {
	networkName := resolveNetworkName(ctx)
	deviceName := resolveDeviceName(ctx)
	pointName := resolvePointName(ctx)
	args := buildPointArgs(ctx)
	q, err := a.DB.GetPointByName(networkName, deviceName, pointName, args)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) PointWriteByNameArgs(ctx *gin.Context) {
	body, _ := getBODYPointWriter(ctx)
	networkName, deviceName, pointName := networkDevicePointNames(ctx)
	q, err := a.DB.PointWriteByName(networkName, deviceName, pointName, body)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) PointWriteByName(ctx *gin.Context) {
	body, _ := getBODYPointWriter(ctx)
	networkName := resolveNetworkName(ctx)
	deviceName := resolveDeviceName(ctx)
	pointName := resolvePointName(ctx)
	q, err := a.DB.PointWriteByName(networkName, deviceName, pointName, body)
	ResponseHandler(q, err, ctx)
}

func networkDevicePointNames(ctx *gin.Context) (networkName, deviceName, pointName string) {
	var aType = ArgsType
	var aDefault = ArgsDefault
	networkName = ctx.DefaultQuery(aType.NetworkName, aDefault.NetworkName)
	deviceName = ctx.DefaultQuery(aType.DeviceName, aDefault.DeviceName)
	pointName = ctx.DefaultQuery(aType.PointName, aDefault.PointName)
	return networkName, deviceName, pointName
}

func (a *PointAPI) CreatePointMetaTags(ctx *gin.Context) {
	pointUUID := resolveID(ctx)
	body, _ := getBodyBulkPointMetaTag(ctx)
	q, err := a.DB.CreatePointMetaTags(pointUUID, body)
	if err != nil {
		ResponseHandler(q, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) DeleteOnePointByArgs(ctx *gin.Context) {
	args := buildDeviceArgs(ctx)
	q, err := a.DB.DeleteOnePointByArgs(args)
	ResponseHandler(q, err, ctx)
}
