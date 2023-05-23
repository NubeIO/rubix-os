package main

import (
	"github.com/NubeIO/lib-schema/systemschema"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/plugin"
	"github.com/NubeIO/rubix-os/plugin/nube/system/smodel"
	"net/http"

	"github.com/gin-gonic/gin"
)

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

const (
	schemaNetwork     = "/schema/network"
	schemaDevice      = "/schema/device"
	schemaPoint       = "/schema/point"
	jsonSchemaNetwork = "/schema/json/network"
	jsonSchemaDevice  = "/schema/json/device"
	jsonSchemaPoint   = "/schema/json/point"
	help              = "/help"
	helpHTML          = "/help/guide"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.GET("/system/schedule/store/:name", func(ctx *gin.Context) {
		obj, ok := inst.store.Get(resolveName(ctx))
		if ok != true {
			ctx.JSON(http.StatusBadRequest, "no schedule exists")
		} else {
			ctx.JSON(http.StatusOK, obj)
		}
	})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		fns, err := inst.db.GetFlowNetworks(api.Args{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		fnsNames := make([]string, 0)
		for _, fn := range fns {
			fnsNames = append(fnsNames, fn.Name)
		}
		networkSchema := smodel.GetNetworkSchema()
		networkSchema.AutoMappingFlowNetworkName.Options = fnsNames
		ctx.JSON(http.StatusOK, networkSchema)
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, smodel.GetDeviceSchema())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, smodel.GetPointSchema())
	})

	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		fns, err := inst.db.GetFlowNetworks(api.Args{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		networkSchema := systemschema.GetNetworkSchema()
		networkSchema.AutoMappingFlowNetworkName.Options = plugin.GetFlowNetworkNames(fns)
		ctx.JSON(http.StatusOK, networkSchema)
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, systemschema.GetDeviceSchema())
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, systemschema.GetPointSchema())
	})
}
