package main

import (
	"fmt"
	pkgmodel "github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
	plgrest "github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/restclient"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getBODYNetwork(ctx *gin.Context) (dto *pkgmodel.Server, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath

	mux.GET("/bacnet/server", func(ctx *gin.Context) {
		cli := plgrest.NewNoAuth(ip, port)
		p, err := cli.GetServer()
		if err != nil {
			fmt.Println(err, "ERROR ON GetServer")
		}
		ctx.JSON(http.StatusOK, p)
	})

	mux.PATCH("/bacnet/server", func(ctx *gin.Context) {
		body, _ := getBODYNetwork(ctx)
		cli := plgrest.NewNoAuth(ip, port)
		p, err := cli.EditServer(*body)
		if err != nil {
			fmt.Println(err, "ERROR ON GET POINTS")
		}
		ctx.JSON(http.StatusOK, p)
	})

	mux.GET("/bacnet/points", func(ctx *gin.Context) {
		cli := plgrest.NewNoAuth(ip, port)
		p, err := cli.GetPoints()
		if err != nil {
			fmt.Println(err, "ERROR ON GET POINTS")
		}
		ctx.JSON(http.StatusOK, p)
	})

	//delete all the bacnet-server points
	mux.DELETE("/bacnet/points", func(ctx *gin.Context) {
		cli := plgrest.NewNoAuth(ip, port)
		p, err := cli.GetPoints()
		for _, pnt := range *p {
			_, err := i.bacnetServerDeletePoint(&pnt)
			if err != nil {
				return
			}
		}
		if err != nil {
			fmt.Println(err, "ERROR ON GET POINTS")
		}
		ctx.JSON(http.StatusOK, p)
	})

}
