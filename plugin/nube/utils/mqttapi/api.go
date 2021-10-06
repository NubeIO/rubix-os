package main

import (
	"github.com/gin-gonic/gin"
)

func resolveToken(ctx *gin.Context) string {
	return ctx.Param("token")
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.POST("/mqtt/download/builds/:token", func(ctx *gin.Context) {

	})

}
