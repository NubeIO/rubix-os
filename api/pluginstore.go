package api

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/services/appstore"
	"github.com/gin-gonic/gin"
)

type PluginStoreApi struct {
	Store *appstore.Store
}

func (a *PluginStoreApi) GetPluginsStorePlugins(c *gin.Context) {
	data, err := a.Store.GetPluginsStorePlugins()
	ResponseHandler(data, err, c)
}

func (a *PluginStoreApi) UploadPluginStorePlugin(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	m := &interfaces.Upload{
		File: file,
	}
	data, err := a.Store.UploadPluginStorePlugin(m)
	ResponseHandler(data, err, c)
}
