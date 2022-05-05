package main

import (
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lwmodel"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lwrest"
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
func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	cli := lwrest.NewChirp(chirpName, chirpPass, ip, port)

	mux.GET("/lorawan/organizations", func(ctx *gin.Context) {
		p, err := cli.GetOrganizations()
		if err != nil {
			log.Info(err, "ERROR ON organizations")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/gateways", func(ctx *gin.Context) {
		p, err := cli.GetGateways()
		if err != nil {
			log.Info(err, "ERROR ON gateways")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/applications", func(ctx *gin.Context) {
		p, err := cli.GetApplications()
		if err != nil {
			log.Info(err, "ERROR ON applications")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/devices", func(ctx *gin.Context) {
		p, err := cli.GetDevices()
		if err != nil {
			log.Info(err, "ERROR ON applications")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		p, err := cli.GetDevice(eui)
		if err != nil {
			log.Info(err, "ERROR ON GetDevice")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.PUT("/lorawan/devices/:eui", func(ctx *gin.Context) {
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

	mux.DELETE("/lorawan/devices/:eui", func(ctx *gin.Context) {
		eui := resolveID(ctx)
		p, err := cli.DeleteDevice(eui)
		if err != nil {
			log.Info(err, "ERROR ON DeleteDevice")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}

	})

	mux.DELETE("/lorawan/devices/drop", func(ctx *gin.Context) {
		_, err := inst.DropDevices()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, "ok")
		}
	})

	mux.GET("/lorawan/device-profiles", func(ctx *gin.Context) {
		p, err := cli.GetDeviceProfiles()
		if err != nil {
			log.Info(err, "ERROR ON device-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/service-profiles", func(ctx *gin.Context) {
		p, err := cli.GetServiceProfiles()
		if err != nil {
			log.Info(err, "ERROR ON service-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.GET("/lorawan/gateway-profiles", func(ctx *gin.Context) {
		p, err := cli.GetGatewayProfiles()
		if err != nil {
			log.Info(err, "ERROR ON gateway-profiles")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})

	mux.PATCH(plugin.PointsWriteURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		uuid := plugin.ResolveID(ctx)
		point, err := inst.writePoint(uuid, body)
		plugin.ResponseHandler(point, err, 0, ctx)
	})

}
