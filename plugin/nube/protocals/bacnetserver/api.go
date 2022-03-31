package main

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

func getBODYNetwork(ctx *gin.Context) (dto *bacnet_model.Server, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPoints(ctx *gin.Context) (dto *model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPointsBacnet(ctx *gin.Context) (dto *bacnet_model.BacnetPoint, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveUUID(ctx *gin.Context) string {
	return ctx.Param("uuid")
}

func resolveObject(ctx *gin.Context) string {
	return ctx.Param("object")
}

func resolveAddress(ctx *gin.Context) string {
	return ctx.Param("address")
}

// RegisterWebhook implements plugin.Webhooker
func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.POST(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.addNetwork(body)
		plugin.ResponseHandler(network, err, 0, ctx)
	})
	mux.POST(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.addDevice(body)
		plugin.ResponseHandler(device, err, 0, ctx)
	})
	mux.POST(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.addPoint(body)
		plugin.ResponseHandler(point, err, 0, ctx)
	})
	//mux.PATCH(plugin.NetworksURL, func(ctx *gin.Context) {
	//	body, _ := plugin.GetBODYNetwork(ctx)
	//	network, err := inst.updateNetwork(body)
	//	plugin.ResponseHandler(network, err, 0, ctx)
	//})
	//mux.PATCH(plugin.DevicesURL, func(ctx *gin.Context) {
	//	body, _ := plugin.GetBODYDevice(ctx)
	//	device, err := inst.updateDevice(body)
	//	plugin.ResponseHandler(device, err, 0, ctx)
	//})
	//mux.PATCH(plugin.PointsURL, func(ctx *gin.Context) {
	//	body, _ := plugin.GetBODYPoint(ctx)
	//	point, err := inst.updatePoint(body)
	//	plugin.ResponseHandler(point, err, 0, ctx)
	//})
	//mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
	//	body, _ := plugin.GetBODYNetwork(ctx)
	//	ok, err := inst.deleteNetwork(body)
	//	plugin.ResponseHandler(ok, err, 0, ctx)
	//})
	//mux.DELETE(plugin.DevicesURL, func(ctx *gin.Context) {
	//	body, _ := plugin.GetBODYDevice(ctx)
	//	ok, err := inst.deleteDevice(body)
	//	plugin.ResponseHandler(ok, err, 0, ctx)
	//})
	//mux.DELETE(plugin.PointsURL, func(ctx *gin.Context) {
	//	body, _ := plugin.GetBODYPoint(ctx)
	//	ok, err := inst.deletePoint(body)
	//	plugin.ResponseHandler(ok, err, 0, ctx)
	//})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, bacnet_model.GetNetworkSchema())
	})

	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, bacnet_model.GetDeviceSchema())
	})

	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, bacnet_model.GetPointSchema())
	})
}
