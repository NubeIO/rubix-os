package main

import (
	"github.com/gin-gonic/gin"
)



// RegisterWebhook implements plugin.Webhooker
func (c *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	c.basePath = basePath
	mux.GET("/message", func(ctx *gin.Context) {

	})
}
