package main

import (
	"fmt"
	plgrest "github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/restclient"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.POST("/bacnet/points", func(ctx *gin.Context) {
		cli := plgrest.NewNoAuth(ip, port)
		p, err := cli.GetPoints()
		if err != nil {
			fmt.Println(err, "ERROR ON GET POINTS")
		}
		ctx.JSON(http.StatusOK, p)
	})

}
