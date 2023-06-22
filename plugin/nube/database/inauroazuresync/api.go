package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/mbmodel"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	help              = "/help"
	schemaNetwork     = "/schema/network"
	jsonSchemaNetwork = "/schema/json/network"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.POST(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.addNetwork(body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.PATCH(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.updateNetwork(body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		ok, err := inst.deleteNetwork(body)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mbmodel.GetNetworkSchema())
	})
	/*
		mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
			fns, err := inst.db.GetFlowNetworks(api.Args{})
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
				return
			}
			networkSchema := modbschema.GetNetworkSchema()
			networkSchema.AutoMappingFlowNetworkName.Options = plugin.GetFlowNetworkNames(fns)
			ctx.JSON(http.StatusOK, networkSchema)
		})

	*/
}
