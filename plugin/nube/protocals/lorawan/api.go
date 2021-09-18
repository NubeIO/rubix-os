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

	mux.GET("/organizations", func(ctx *gin.Context) {
		p, err := cli.GetOrganizations()
		if err != nil {
			log.Info(err, "ERROR ON organizations")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/gateways", func(ctx *gin.Context) {
		p, err := cli.GetGateways()
		if err != nil {
			log.Info(err, "ERROR ON gateways")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/applications", func(ctx *gin.Context) {
		p, err := cli.GetApplications()
		if err != nil {
			log.Info(err, "ERROR ON applications")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/devices", func(ctx *gin.Context) {
		p, err := cli.GetDevices()
		if err != nil {
			log.Info(err, "ERROR ON applications")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		p, err := cli.GetDevice(eui)
		if err != nil {
			log.Info(err, "ERROR ON GetDevice")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.PUT("/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		device, err := bodyDevice(ctx)
		if err != nil {
			return
		}
		_, err = cli.EditDevice(eui, device)
		if err != nil {
			log.Info(err, "ERROR ON GetDevice")
			ctx.JSON(http.StatusBadRequest, "fail")
		} else {
			ctx.JSON(http.StatusOK, "ok")
		}
	})

	mux.DELETE("/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		p, err := cli.DeleteDevice(eui)
		if err != nil {
			log.Info(err, "ERROR ON DeleteDevice")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}

	})

	mux.DELETE("/devices/drop", func(ctx *gin.Context) {
		_, err := i.DropDevices()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, "ok")
		}
	})

	mux.GET("/device-profiles", func(ctx *gin.Context) {
		p, err := cli.GetDeviceProfiles()
		if err != nil {
			log.Info(err, "ERROR ON device-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/service-profiles", func(ctx *gin.Context) {
		p, err := cli.GetServiceProfiles()
		if err != nil {
			log.Info(err, "ERROR ON service-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/gateway-profiles", func(ctx *gin.Context) {
		p, err := cli.GetGatewayProfiles()
		if err != nil {
			log.Info(err, "ERROR ON gateway-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

}
