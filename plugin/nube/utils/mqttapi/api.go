package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func resolveToken(ctx *gin.Context) string {
	return ctx.Param("token")
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.POST("/mqtt/download/builds/:token", func(ctx *gin.Context) {
		token := resolveToken(ctx)
		err := download(token)
		if err != nil {
			log.Error(err, "ERROR ON download builds")
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, "downloaded ok!")
		}

	})

}
