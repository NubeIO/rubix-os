package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	help         = "/lora/help"
	listSerial   = "/lora/list/serial"
	wizardSerial = "/lora/wizard/serial"
)

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
	mux.GET(listSerial, func(ctx *gin.Context) {
		serial, err := i.listSerialPorts()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, serial)
		}
	})
	mux.POST(wizardSerial, func(ctx *gin.Context) {
		serial, err := i.wizardSerial()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {

			ctx.JSON(http.StatusOK, serial)
		}
	})

}
