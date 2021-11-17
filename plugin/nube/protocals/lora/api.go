package main

import (
	serial_model "github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/model"
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

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.GET(help, func(ctx *gin.Context) {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, "add")
		}
	})
	mux.POST(restartSerial, func(ctx *gin.Context) {
		err := i.SerialClose()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			err := i.SerialOpen()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				ctx.JSON(http.StatusOK, "ok")
			}
		}
	})
	mux.GET(listSerial, func(ctx *gin.Context) {
		serial, err := i.listSerialPorts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST(wizardSerial, func(ctx *gin.Context) {
		body, err := bodyWizard(ctx)
		serial, err := i.wizardSerial(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, serial_model.GetNetworkSchema())
	})

	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, serial_model.GetDeviceSchema())
	})

	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, serial_model.GetPointSchema())
	})
}
