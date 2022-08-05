package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	listBuckets        = "/backup/buckets/list"
	backUpdateNetworks = "/backup/networks"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.GET(listBuckets, func(ctx *gin.Context) {
		err := inst.connection()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			buckets, err := inst.minioClient.ListBuckets()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			}
			ctx.JSON(http.StatusOK, buckets)
		}
	})
	mux.POST(backUpdateNetworks, func(ctx *gin.Context) {
		err := inst.connection()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			err := inst.backNetworks()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err)
			} else {
				ctx.JSON(http.StatusOK, "upload ok")
			}
		}
	})
}
