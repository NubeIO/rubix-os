package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/plgrest"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nrest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube_api"
	nube_api_bacnetserver "github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube_api/bacnetserver"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
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

func resolveUUID(ctx *gin.Context) string {
	return ctx.Param("uuid")
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
	mux.POST(plugin.NetworksURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYNetwork(ctx)
		network, err := inst.addNetwork(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, network)
		}
	})
	mux.DELETE(plugin.NetworksURL, func(ctx *gin.Context) {
		network, err := inst.deleteNetwork()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, network)
		}
	})
	mux.DELETE(plugin.DevicesURL, func(ctx *gin.Context) {
		ok, err := inst.deleteDevice()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, ok)
		}
	})
	mux.DELETE(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		ok, err := inst.deletePoint(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, ok)
		}
	})
	mux.POST(plugin.DevicesURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYDevice(ctx)
		device, err := inst.addDevice(body)
		fmt.Println(err)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, device)
		}
	})
	mux.POST(plugin.PointsURL, func(ctx *gin.Context) {
		body, _ := plugin.GetBODYPoint(ctx)
		p, err := inst.addPoint(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
		} else {
			ctx.JSON(http.StatusOK, p)
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

		//inc generic reset client
		rc := &nrest.ReqType{
			BaseUri: "0.0.0.0",
			LogPath: "helpers.nrest",
		}

		//inc nube rest client
		c := &nube_api.NubeRest{
			Rest:          rc,
			RubixPort:     nube_api.DefaultRubixService,
			RubixUsername: "admin",
			RubixPassword: "N00BWires",
			UseRubixProxy: true,
		}

		//new nube rest client
		nubeRest := nube_api.New(c)
		//nubeRest.GetToken()

		//bacnet client
		options := &nrest.ReqOpt{
			Timeout:          500 * time.Second,
			RetryCount:       1,
			RetryWaitTime:    1 * time.Second,
			RetryMaxWaitTime: 0,
			Headers:          map[string]interface{}{"Authorization": nubeRest.RubixToken},
		}
		rc.Service = "bacnet-server"
		rc.LogPath = "helpers.nrest.bacnet.server"
		rc.Port = nube_api.DefaultPortBacnet
		c.RubixProxyPath = nube_api.ProxyBacnet
		bacnetClient := &nube_api_bacnetserver.RestClient{
			NubeRest: nubeRest,
			Options:  options,
		}
		//get points
		//_, r := bacnetClient.GetPoint("BhLtrFaNrtBxhVLyjc5CHi")
		//_, r := bacnetClient.GetPoints()

		body, _ := getBODYPoints(ctx)
		var bacPoint nube_api_bacnetserver.BacnetPoint
		if body.Description == "" {
			bacPoint.Description = "na"
		}
		bacPoint.ObjectName = body.Name
		bacPoint.Enable = true
		bacPoint.Address = utils.IntIsNil(body.AddressID)
		bacPoint.ObjectType = body.ObjectType
		bacPoint.COV = utils.Float64IsNil(body.COV)
		bacPoint.EventState = "normal"
		bacPoint.Units = "noUnits"
		bacPoint.RelinquishDefault = utils.Float64IsNil(body.Fallback)
		_, r := bacnetClient.AddPoint(bacPoint)

		if r.ApiReply.Err != nil {
			ctx.JSON(r.Response.StatusCode, r.Response)
		} else {
			ctx.JSON(r.Response.StatusCode, r.Response)
		}

	})
	mux.POST("/bacnet/points", func(ctx *gin.Context) {
		body, _ := getBODYPoints(ctx)
		var bacPoint nube_api_bacnetserver.BacnetPoint
		if body.Description == "" {
			bacPoint.Description = "na"
		}
		bacPoint.ObjectName = body.Name
		bacPoint.Enable = true
		bacPoint.Address = utils.IntIsNil(body.AddressID)
		bacPoint.ObjectType = body.ObjectType
		bacPoint.COV = utils.Float64IsNil(body.COV)
		bacPoint.EventState = "normal"
		bacPoint.Units = "noUnits"
		bacPoint.RelinquishDefault = utils.Float64IsNil(body.Fallback)
		_, r := bacnetClient.AddPoint(bacPoint)
		if r.ApiReply.Err != nil {
			ctx.JSON(r.Response.StatusCode, r.Response)
		} else {

			ctx.JSON(r.Response.StatusCode, r.Response)
		}
	})
	mux.PATCH("/bacnet/points/:uuid", func(ctx *gin.Context) {
		uuid := resolveUUID(ctx)
		body, _ := getBODYPoints(ctx)
		point, err := inst.pointPatch2(body)
		if err != nil {
			fmt.Println(111111, err)
		}
		_, r := bacnetClient.UpdatePoint(uuid, point)
		if r.ApiReply.Err != nil {
			ctx.JSON(r.Response.StatusCode, r.Response)
		} else {
			ctx.JSON(r.Response.StatusCode, r.Response)
		}
	})

	mux.DELETE("/bacnet/points/:uuid", func(ctx *gin.Context) {
		uuid := resolveUUID(ctx)
		r := bacnetClient.DeletePoint(uuid)
		if r.ApiReply.Err != nil {
			ctx.JSON(r.Response.StatusCode, r.Response)
		} else {
			ctx.JSON(r.Response.StatusCode, r.Response)
			return
		}
	})

	//delete all the bacnet-server points
	mux.DELETE("/bacnet/points/drop", func(ctx *gin.Context) {

		r := bacnetClient.DropPoints()
		if r.ApiReply.Err != nil {
			ctx.JSON(r.Response.StatusCode, r.Response)
		} else {
			ctx.JSON(r.Response.StatusCode, r.Response)
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
