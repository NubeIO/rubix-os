package main

import (
	rest "github.com/NubeDev/flow-framework/plugin/nube/protocals/lorawan/restclient"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// EUI body of device eui
type EUI struct {
	EUI string `json:"eui"`
}

func bodyEUI(ctx *gin.Context) (dto *EUI, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

const chirpName = "admin"
const chirpPass = "admin"

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	cli := rest.NewChirp(chirpName, chirpPass, ip, port)

	mux.GET("/organizations", func(ctx *gin.Context) {
		p, err := cli.GetOrganizations()
		if err != nil {
			log.Info(err, "ERROR ON organizations")
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.JSON(http.StatusOK, p)
	})

	mux.GET("/gateways", func(ctx *gin.Context) {
		p, err := cli.GetGateways()
		if err != nil {
			log.Info(err, "ERROR ON gateways")
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.JSON(http.StatusOK, p)
	})

	mux.GET("/applications", func(ctx *gin.Context) {
		p, err := cli.GetApplications()
		if err != nil {
			log.Info(err, "ERROR ON applications")
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.JSON(http.StatusOK, p)
	})

	mux.GET("/devices", func(ctx *gin.Context) {
		p, err := cli.GetDevices()
		if err != nil {
			log.Info(err, "ERROR ON applications")
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.JSON(http.StatusOK, p)
	})

	mux.DELETE("/devices", func(ctx *gin.Context) {
		body, _ := bodyEUI(ctx)
		p, err := cli.DeleteDevice(body.EUI)
		if err != nil {
			log.Info(err, "ERROR ON applications")
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.JSON(http.StatusOK, p)
	})

	mux.DELETE("/devices/drop", func(ctx *gin.Context) {
		_, err := i.DropDevices()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.String(http.StatusOK, "ok")
	})

	mux.GET("/device-profiles", func(ctx *gin.Context) {
		p, err := cli.GetDeviceProfiles()
		if err != nil {
			log.Info(err, "ERROR ON device-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.JSON(http.StatusOK, p)
	})

	mux.GET("/service-profiles", func(ctx *gin.Context) {
		p, err := cli.GetServiceProfiles()
		if err != nil {
			log.Info(err, "ERROR ON service-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.JSON(http.StatusOK, p)
	})

	mux.GET("/gateway-profiles", func(ctx *gin.Context) {
		p, err := cli.GetGatewayProfiles()
		if err != nil {
			log.Info(err, "ERROR ON gateway-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		}
		ctx.JSON(http.StatusOK, p)
	})

}
