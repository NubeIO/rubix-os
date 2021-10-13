package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	listBuckets        = "/backup/buckets/list"
	backUpdateNetworks = "/backup/networks"
)

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.GET(listBuckets, func(ctx *gin.Context) {
		err := i.connection()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			buckets, err := i.minioClient.ListBuckets()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			}
			ctx.JSON(http.StatusOK, buckets)
		}
	})
	mux.POST(backUpdateNetworks, func(ctx *gin.Context) {
		err := i.connection()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			err := i.backNetworks()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				ctx.JSON(http.StatusOK, "upload ok")
			}
		}
	})
}
