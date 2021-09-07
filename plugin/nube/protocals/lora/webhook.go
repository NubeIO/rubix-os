package main

import (
	"github.com/gin-gonic/gin"
)


type message struct {
	Message  string                 `json:"message" query:"message" form:"message"`
	Title    string                 `json:"title" query:"title" form:"title"`
	Priority int                    `json:"priority" query:"priority" form:"priority"`
	Extras   map[string]interface{} `json:"extras" query:"-" form:"-"`
}




// RegisterWebhook implements plugin.Webhooker
func (c *PluginTest) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	c.basePath = basePath
	mux.GET("/message", func(ctx *gin.Context) {

	})
}
