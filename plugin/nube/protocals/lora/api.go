package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)




// RegisterWebhook implements plugin.Webhooker
func (c *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	c.basePath = basePath
	//restart plugin
	mux.POST("/restart", func(ctx *gin.Context) {
		err := c.Disable();if err != nil {
			log.Error("LORA: error on restart (disable) plugin %s", err)
			ctx.JSON(http.StatusBadRequest, "restart fail")
		}
		time.Sleep(300 * time.Millisecond)
		err = c.Enable();if err != nil {
			log.Error("LORA: error on restart (enable) plugin %s", err)
			ctx.JSON(http.StatusBadRequest, "restart fail")
		}
		ctx.JSON(http.StatusOK, "restart ok")
	})

}
