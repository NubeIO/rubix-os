package main

import (
	"fmt"
	"net/http"

	"github.com/NubeIO/lib-schema/lorawanschema"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/plugin"
	"github.com/NubeIO/rubix-os/plugin/nube/protocals/lorawan/csmodel"
	"github.com/NubeIO/rubix-os/plugin/nube/protocals/lorawan/csrest"
	"github.com/gin-gonic/gin"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.POST(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		err := inst.addDevice(body)
		api.ResponseHandler(body, err, ctx)
	})
	mux.PATCH(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.db.UpdateNetwork(body.UUID, body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.PATCH(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		err := inst.updateDevice(body)
		api.ResponseHandler(body, err, ctx)
	})
	mux.PATCH(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.db.UpdatePoint(body.UUID, body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
		inst.deleteNetwork()
		api.ResponseHandler(true, nil, ctx)
	})
	mux.DELETE(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		err := inst.deleteDevice(body)
		api.ResponseHandler(true, err, ctx)
	})
	mux.DELETE(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		ok, err := inst.db.DeletePoint(body.UUID)
		api.ResponseHandler(ok, err, ctx)
	})

	mux.GET(plugin.SchemaLegacyNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, csmodel.GetNetworkSchema())
	})
	mux.GET(plugin.SchemaLegacyDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, csmodel.GetDeviceSchema())
	})
	mux.GET(plugin.SchemaLegacyPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, csmodel.GetPointSchema())
	})

	mux.GET(plugin.JsonSchemaNetwork, func(ctx *gin.Context) {
		fns, err := inst.db.GetFlowNetworks(api.Args{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		networkSchema := lorawanschema.GetNetworkSchema()
		networkSchema.AutoMappingFlowNetworkName.Options = plugin.GetFlowNetworkNames(fns)
		ctx.JSON(http.StatusOK, networkSchema)
	})
	mux.GET(plugin.JsonSchemaDevice, func(ctx *gin.Context) {
		schema := lorawanschema.GetDeviceSchema()
		inst.fillDeviceProfilesSchema(schema)
		ctx.JSON(http.StatusOK, schema)
	})
	mux.GET(plugin.JsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, lorawanschema.GetPointSchema())
	})

	// ---------------------- CS APIs ------------------------------

	mux.GET(fmt.Sprintf("%s/*any", csrest.CsURLPrefix), inst.chirpStack.Proxy)
	mux.POST(fmt.Sprintf("%s/*any", csrest.CsURLPrefix), inst.chirpStack.Proxy)
	mux.PUT(fmt.Sprintf("%s/*any", csrest.CsURLPrefix), inst.chirpStack.Proxy)
	mux.PATCH(fmt.Sprintf("%s/*any", csrest.CsURLPrefix), inst.chirpStack.Proxy)
	mux.DELETE(fmt.Sprintf("%s/*any", csrest.CsURLPrefix), inst.chirpStack.Proxy)
}
