package main

import (
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
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
		msg := new(message)
		if err := ctx.Bind(msg); err == nil {
			err := c.msgHandler.SendMessage(pluginapi.Message{
				Message:  msg.Message,
				Title:    msg.Title,
				Priority: msg.Priority,
				Extras:   msg.Extras,
			})
			if err != nil {
				return
			}
			ctx.JSON(200, msg)
		}
	})
}
