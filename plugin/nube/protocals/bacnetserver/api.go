package main

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/plgrest"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

func getBODYNetwork(ctx *gin.Context) (dto *bacnet_model.Server, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPoints(ctx *gin.Context) (dto *model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPointsBacnet(ctx *gin.Context) (dto *bacnet_model.BacnetPoint, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveObject(ctx *gin.Context) string {
	return ctx.Param("object")
}

func resolveAddress(ctx *gin.Context) string {
	return ctx.Param("address")
}

// RegisterWebhook implements plugin.Webhooker
func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.POST("/points", func(ctx *gin.Context) {
		body, _ := getBODYPoints(ctx)
		point, httpRes, err := inst.addPoint(body)
		if httpRes != nil {
			ctx.JSON(httpRes.StatusCode, httpRes.AsJsonNoErr())
			return
		}
		if err != nil {
			log.Error(err, "ERROR ON PingServer")
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, point)
		}
	})

	mux.GET("/bacnet/ping", func(ctx *gin.Context) {
		cli := plgrest.NewNoAuth(ip, string(port))
		p, err := cli.PingServer()
		if err != nil {
			log.Error(err, "ERROR ON PingServer")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.GET("/bacnet/server", func(ctx *gin.Context) {
		cli := plgrest.NewNoAuth(ip, string(port))
		p, err := cli.GetServer()
		if err != nil {
			log.Error(err, "ERROR ON GetServer")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.PATCH("/bacnet/server", func(ctx *gin.Context) {
		body, _ := getBODYNetwork(ctx)
		cli := plgrest.NewNoAuth(ip, string(port))
		p, err := cli.EditServer(*body)
		if err != nil {
			log.Error(err, "ERROR ON EditServer")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	//POINTS
	mux.GET("/bacnet/points", func(ctx *gin.Context) {
		//points, res := inst.restClient.GetPoints()
		//statusCode := res.Reply.StatusCode
		//if res.IsError {
		//	ctx.JSON(statusCode, res.Reply.Err)
		//	return
		//} else if res.IsError {
		//	//ctx.JSON(statusCode, res.Response)
		//	return
		//} else {
		//	ctx.JSON(statusCode, points)
		//}

	})
	mux.POST("/bacnet/points", func(ctx *gin.Context) {
		body, _ := getBODYPointsBacnet(ctx)
		cli := plgrest.NewNoAuth(ip, string(port))
		p, err := cli.AddPoint(*body)
		if err != nil {
			log.Error(err, "ERROR ON AddPoint")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.PATCH("/bacnet/points/:object/:address", func(ctx *gin.Context) {
		body, _ := getBODYPointsBacnet(ctx)
		obj := resolveObject(ctx)
		addr := resolveAddress(ctx)
		cli := plgrest.NewNoAuth(ip, string(port))
		p, err := cli.EditPoint(*body, obj, utils.ToInt(addr))
		if err != nil {
			log.Error(err, "ERROR ON EditPoint")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	//delete all the bacnet-server points
	mux.DELETE("/bacnet/points/drop", func(ctx *gin.Context) {
		cli := plgrest.NewNoAuth(ip, string(port))
		p, err := cli.GetPoints()
		for _, pnt := range *p {
			_, err := inst.bacnetServerDeletePoint(&pnt)
			if err != nil {
				return
			}
		}
		if err != nil {
			log.Error(err, "ERROR ON bacnetServerDeletePoint")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.POST("/bacnet/wizard", func(ctx *gin.Context) {
		wizard, err := inst.wizard()
		if err != nil {
			log.Error(err, "ERROR ON wizard")
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, wizard)
			return
		}
	})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, bacnet_model.GetNetworkSchema())
	})

	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, bacnet_model.GetDeviceSchema())
	})

	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, bacnet_model.GetPointSchema())
	})
}
