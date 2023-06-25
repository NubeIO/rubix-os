package main

import (
	"github.com/gin-gonic/gin"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	sites := mux.Group("/sites")

	sites.POST("", inst.CreateSite)
	sites.GET("", inst.GetAllSites)
	sites.GET("/:id", inst.GetSite)
	sites.POST("/name", inst.GetSiteByName)
	sites.POST("/address", inst.GetSiteByAddress)
	sites.PATCH("/:id", inst.UpdateSite)
	sites.DELETE("/:id", inst.DeleteSite)
}
