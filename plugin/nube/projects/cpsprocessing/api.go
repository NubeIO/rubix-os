package main

import (
	"github.com/gin-gonic/gin"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	sites := mux.Group("/sites")

	sites.POST("", inst.CreateSite)
	sites.GET("", inst.GetAllSites)
	sites.GET("/:site_ref", inst.GetSite)
	sites.POST("/name", inst.GetSiteByName)
	sites.POST("/address", inst.GetSiteByAddress)
	sites.PATCH("/:site_ref", inst.UpdateSite)
	sites.DELETE("/:site_ref", inst.DeleteSite)

	thresholds := mux.Group("/thresholds")

	thresholds.POST("", inst.CreateThreshold)
	thresholds.GET("/:site_ref", inst.GetLastThresholdBySiteRef)
}
