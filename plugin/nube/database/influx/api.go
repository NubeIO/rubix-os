package main

import (
	lwmodel "github.com/NubeDev/flow-framework/plugin/nube/protocals/lorawan/model"
	rest "github.com/NubeDev/flow-framework/plugin/nube/protocals/lorawan/restclient"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func bodyDevice(ctx *gin.Context) (dto lwmodel.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveID(ctx *gin.Context) string {
	return ctx.Param("eui")
}

const chirpName = "admin"
const chirpPass = "admin"

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	cli := rest.NewChirp(chirpName, chirpPass, ip, port)

	mux.GET("/influx/temperatures", func(ctx *gin.Context) {
		p, err := cli.GetOrganizations()
		if err != nil {
			log.Info(err, "ERROR ON organizations")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.POST("/influx/temperature", func(ctx *gin.Context) {
		p, err := cli.GetOrganizations()
		if err != nil {
			log.Info(err, "ERROR ON organizations")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

}
