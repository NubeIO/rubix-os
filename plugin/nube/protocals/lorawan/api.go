package main

import (
	"net/http"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	"github.com/gin-gonic/gin"
)

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

// RegisterWebhook implements plugin.Webhooker
func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, csmodel.GetNetworkSchema())
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, csmodel.GetDeviceSchema())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, csmodel.GetPointSchema())
	})

	mux.PATCH(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.db.UpdateNetwork(body.UUID, body, true)
		api.ResponseHandler(network, err, ctx)
	})
	mux.PATCH(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.db.UpdateDevice(body.UUID, body, true)
		api.ResponseHandler(device, err, ctx)
	})
	mux.PATCH(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.db.UpdatePoint(body.UUID, body, true)
		api.ResponseHandler(point, err, ctx)
	})
	mux.DELETE(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		ok, err := inst.db.DeletePoint(body.UUID)
		api.ResponseHandler(ok, err, ctx)
	})
}
