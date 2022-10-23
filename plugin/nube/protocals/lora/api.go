package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lora/loramodel"
	"github.com/NubeIO/lib-schema/loraschema"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	help          = "/help"
	restartSerial = "/serial/restart"
	listSerial    = "/serial/list"
	wizardSerial  = "/wizard/serial"
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

const (
	jsonSchemaNetwork = "/schema/json/network"
	jsonSchemaDevice  = "/schema/json/device"
	jsonSchemaPoint   = "/schema/json/point"
)

// wizard
type wizard struct {
	SerialPort string `json:"serial_port"`
	SensorID   string `json:"sensor_id"`
	SensorType string `json:"sensor_type"`
}

func bodyWizard(ctx *gin.Context) (dto wizard, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

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
		network, err := inst.db.UpdateNetwork(body.UUID, body, true)
		api.ResponseHandler(network, err, ctx)
	})
	mux.PATCH(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.db.UpdateDevice(body.UUID, body, true)
		inst.updateDevicePointsAddress(device)
		api.ResponseHandler(device, err, ctx)
	})
	mux.PATCH(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		point, err := inst.db.UpdatePoint(body.UUID, body, true, false)
		api.ResponseHandler(point, err, ctx)
	})
	mux.PATCH(plugin.PointsWriteURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBodyPointWriter(ctx)
		uuid := plugin.ResolveID(ctx)
		point, _, _, _, err := inst.db.PointWrite(uuid, body, false)
		api.ResponseHandler(point, err, ctx)
	})
	mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		ok, err := inst.db.DeleteNetwork(body.UUID)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		ok, err := inst.db.DeleteDevice(body.UUID)
		api.ResponseHandler(ok, err, ctx)
	})
	mux.DELETE(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		ok, err := inst.db.DeletePoint(body.UUID)
		api.ResponseHandler(ok, err, ctx)
	})

	mux.POST(restartSerial, func(ctx *gin.Context) {
		err := inst.Disable()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			err := inst.Enable()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				ctx.JSON(http.StatusOK, "ok")
			}
		}
	})
	mux.GET(listSerial, func(ctx *gin.Context) {
		serial, err := inst.listSerialPorts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST(wizardSerial, func(ctx *gin.Context) {
		body, _ := bodyWizard(ctx)
		serial, err := inst.wizardSerial(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, loramodel.GetNetworkSchema())
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, loramodel.GetDeviceSchema())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, loramodel.GetPointSchema())
	})
	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, loraschema.GetNetworkSchema())
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, loraschema.GetDeviceSchema())
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, loraschema.GetPointSchema())
	})
}
