package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/edgemodel"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/edgerest"
	"github.com/NubeIO/lib-schema/edge28schema"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func bodyWizard(ctx *gin.Context) (dto wizard, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

const (
	schemaNetwork     = "/schema/network"
	schemaDevice      = "/schema/device"
	schemaPoint       = "/schema/point"
	jsonSchemaNetwork = "/schema/json/network"
	jsonSchemaDevice  = "/schema/json/device"
	jsonSchemaPoint   = "/schema/json/point"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.POST(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.addNetwork(body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.POST(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.addDevice(body)
		api.ResponseHandler(device, err, ctx)
	})
	mux.POST(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.addPoint(body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.PATCH(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.updateNetwork(body)
		api.ResponseHandler(network, err, ctx)
	})
	mux.PATCH(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.updateDevice(body)
		api.ResponseHandler(device, err, ctx)
	})
	mux.PATCH(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.updatePoint(body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.PATCH(plugin.PointsWriteURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBodyPointWriter(ctx)
		uuid := plugin.ResolveID(ctx)
		point, err := inst.writePoint(uuid, body)
		api.ResponseHandler(point, err, ctx)
	})
	mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		ok, err := inst.deleteNetwork(body)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		ok, err := inst.deleteDevice(body)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		ok, err := inst.deletePoint(body)
		api.ResponseHandler(ok, err, ctx)
	})

	mux.GET("/edge/ping", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		p, err := rest.PingServer()
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, p)
			return
		}
	})
	mux.GET("/edge/read/ui/all", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		p, err := rest.GetUIs()
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, p)
			return
		}
	})
	mux.GET("/edge/read/di/all", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		p, err := rest.GetDIs()
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.POST("/edge/wizard", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		p, err := inst.wizard(body)
		if err != nil {
			log.Info(err, "ERROR ON ping server")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.POST("/edge/write/uo", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		ioNum := body.IONum
		val := body.Value
		p, err := rest.WriteUO(ioNum, val)
		if err != nil {
			log.Info(err, "ERROR ON write uo")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.POST("/edge/write/do", func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		rest := edgerest.NewNoAuth(body.IP, body.Port)
		ioNum := body.IONum
		val := body.Value
		p, err := rest.WriteDO(ioNum, val)
		if err != nil {
			log.Info(err, "ERROR ON write do")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, edgemodel.GetNetworkSchema())
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, edgemodel.GetDeviceSchema())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, edgemodel.GetPointSchema())
	})
	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		fns, err := inst.db.GetFlowNetworks(api.Args{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		networkSchema := edge28schema.GetNetworkSchema()
		networkSchema.AutoMappingFlowNetworkName.Options = plugin.GetFlowNetworkNames(fns)
		ctx.JSON(http.StatusOK, networkSchema)
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, edge28schema.GetDeviceSchema())
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, edge28schema.GetPointSchema())
	})
}
