package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/gin-gonic/gin"
)

type PointDatabase interface {
	GetPoints(args Args) ([]*model.Point, error)
	GetPointsBulk(bulkPoints []*model.Point) ([]*model.Point, error)
	GetPointsBulkUUIs() ([]string, error)
	GetPoint(uuid string, args Args) (*model.Point, error)
	CreatePoint(body *model.Point) (*model.Point, error)
	UpdatePoint(uuid string, body *model.Point) (*model.Point, error)
	PointWrite(uuid string, body *model.PointWriter) (*model.Point, bool, bool, bool, error)
	GetOnePointByArgs(args Args) (*model.Point, error)
	DeletePoint(uuid string) (bool, error)
	GetPointByName(networkName, deviceName, pointName string, args Args) (*model.Point, error)
	PointWriteByName(networkName, deviceName, pointName string, body *model.PointWriter) (*model.Point, error)
	DeleteOnePointByArgs(args Args) (bool, error)
	DeletePointByName(networkName, deviceName, pointName string, args Args) (bool, error)
	GetPointWithParent(uuid string) (*interfaces.PointWithParent, error)

	CreatePointPlugin(body *model.Point) (*model.Point, error)
	UpdatePointPlugin(uuid string, body *model.Point) (*model.Point, error)
	WritePointPlugin(uuid string, body *model.PointWriter) (*model.Point, error)
	DeletePointPlugin(uuid string) (bool, error)

	CreatePointMetaTags(pointUUID string, pointMetaTags []*model.PointMetaTag) ([]*model.PointMetaTag, error)

	ResolveHost(uuid string, name string) (*model.Host, error)
}
type PointAPI struct {
	DB PointDatabase
}

func (a *PointAPI) GetPoints(ctx *gin.Context) {
	args := buildPointArgs(ctx)
	q, err := a.DB.GetPoints(args)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointsBulkUUIs(ctx *gin.Context) {
	q, err := a.DB.GetPointsBulkUUIs()
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

func (a *PointAPI) DeletePointByName(ctx *gin.Context) {
	networkName := resolveNetworkName(ctx)
	deviceName := resolveDeviceName(ctx)
	pointName := resolvePointName(ctx)
	args := buildPointArgs(ctx)
	q, err := a.DB.DeletePointByName(networkName, deviceName, pointName, args)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointWithParent(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetPointWithParent(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) GetPointByHost(ctx *gin.Context) {
	cli, err := a.resolveClient(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	uuid := resolveID(ctx)
	q, err := cli.GetPoint(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) WritePointByHost(ctx *gin.Context) {
	cli, err := a.resolveClient(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	uuid := resolveID(ctx)
	body, _ := getBODYPointWriter(ctx)
	q, err := cli.WritePoint(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *PointAPI) resolveClient(ctx *gin.Context) (*client.FlowClient, error) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		return nil, err
	}
	return client.NewClient(host.IP, host.Port, host.ExternalToken), nil
}
