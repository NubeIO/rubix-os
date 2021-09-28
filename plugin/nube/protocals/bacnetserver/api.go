package main

import (
	pkgmodel "github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
	plgrest "github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/restclient"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getBODYNetwork(ctx *gin.Context) (dto *pkgmodel.Server, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPoints(ctx *gin.Context) (dto *pkgmodel.BacnetPoint, err error) {
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
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.GET("/bacnet/ping", func(ctx *gin.Context) {
		cli := plgrest.NewNoAuth(ip, port)
		p, err := cli.PingServer()
		if err != nil {
			log.Error(err, "ERROR ON PingServer")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.GET("/bacnet/server", func(ctx *gin.Context) {
		cli := plgrest.NewNoAuth(ip, port)
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
		cli := plgrest.NewNoAuth(ip, port)
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
		cli := plgrest.NewNoAuth(ip, port)
		p, err := cli.GetPoints()
		if err != nil {
			log.Error(err, "ERROR ON GetPoints")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.POST("/bacnet/points", func(ctx *gin.Context) {
		body, _ := getBODYPoints(ctx)
		cli := plgrest.NewNoAuth(ip, port)
		p, err := cli.AddPoint(*body)
		if err != nil {
			log.Error(err, "ERROR ON AddPoint")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, p)
		}
	})
	mux.PATCH("/bacnet/points/:object/:address", func(ctx *gin.Context) {
		body, _ := getBODYPoints(ctx)
		obj := resolveObject(ctx)
		addr := resolveAddress(ctx)
		cli := plgrest.NewNoAuth(ip, port)
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
		cli := plgrest.NewNoAuth(ip, port)
		p, err := cli.GetPoints()
		for _, pnt := range *p {
			_, err := i.bacnetServerDeletePoint(&pnt)
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
		wizard, err := i.wizard()
		if err != nil {
			log.Error(err, "ERROR ON wizard")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, wizard)
		}
	})
}
