package main

import (
	"github.com/NubeIO/rubix-os/schema/systemschema"
	"net/http"

	"github.com/gin-gonic/gin"
)

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

const (
	jsonSchemaNetwork = "/schema/json/network"
	jsonSchemaDevice  = "/schema/json/device"
	jsonSchemaPoint   = "/schema/json/point"
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
	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, systemschema.GetNetworkSchema())
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, systemschema.GetDeviceSchema())
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, systemschema.GetPointSchema())
	})
}
