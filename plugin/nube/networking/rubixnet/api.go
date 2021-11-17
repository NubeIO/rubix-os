package main

import (
	"github.com/NubeIO/flow-framework/src/system/networking"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	localIP    = "/rubixnet/local/ip"
	externalIP = "/rubixnet/external/ip"
)

type MqttPayload struct {
	Value    *float64
	Priority int
}

func getBODYPoints(ctx *gin.Context) (dto *MqttPayload, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.GET(localIP, func(ctx *gin.Context) {
		address, _, _, err := networking.IpAddresses()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, address)
		}
	})
	mux.GET(externalIP, func(ctx *gin.Context) {
		address, err := networking.ExternalIPV4()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, address)
		}
	})
}
