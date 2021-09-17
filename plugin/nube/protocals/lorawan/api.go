package main

import (
	"fmt"
	csrest "github.com/NubeDev/flow-framework/plugin/nube/protocals/lorawan/restclient"

	"github.com/gin-gonic/gin"
	"net/http"
)

//func body(ctx *gin.Context) (dto *csrest.Server, err error) {
//	err = ctx.ShouldBindJSON(&dto)
//	return dto, err
//}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath

	mux.GET("/bacnet/server", func(ctx *gin.Context) {
		cli := csrest.NewChirp("", "", ip, port)
		p, err := cli.GetApplications()
		if err != nil {
			fmt.Println(err, "ERROR ON GetServer")
		}
		ctx.JSON(http.StatusOK, p)
	})

}
