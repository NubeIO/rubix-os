package api

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/src/cli/cligetter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type EdgePluginDatabase interface {
	ResolveHost(uuid string, name string) (*model.Host, error)
}

type EdgePluginApi struct {
	DB EdgeAppDatabase
}

func (a *EdgePluginApi) EdgeListPlugins(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	plugins, connectionErr, requestErr := cli.ListPlugins()
	if connectionErr != nil {
		ctx.JSON(502, interfaces.Message{Message: connectionErr.Error()})
		return
	}
	ResponseHandler(plugins, requestErr, ctx)
}

func (a *EdgePluginApi) EdgeUploadPlugin(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	var m *interfaces.Plugin
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	data, err := cli.PluginUpload(m)
	ResponseHandler(data, err, ctx)
}

func (a *EdgePluginApi) EdgeMoveFromDownloadToInstallPlugins(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	resp, err := cli.MovePluginsFromDownloadToInstallDir()
	ResponseHandler(resp, err, ctx)
}

func (a *EdgePluginApi) EdgeDeletePlugin(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	pluginName := ctx.Param("plugin_name")
	arch := ctx.Query("arch")
	cli := cligetter.GetEdgeBiosClient(host)
	installPluginFilePath := global.Installer.GetAppPluginInstallFilePath(pluginName, arch)
	_, connectionErr, requestErr := cli.DeleteFiles(installPluginFilePath)
	if connectionErr != nil {
		log.Errorf(connectionErr.Error())
		ctx.JSON(502, interfaces.Message{Message: connectionErr.Error()})
		return
	}
	if requestErr != nil {
		ResponseHandler(nil, requestErr, ctx)
		log.Errorf(requestErr.Error())
		return
	}
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("successfully deleted plugin %s", pluginName)}, nil, ctx)
}

func (a *EdgePluginApi) EdgeDeleteDownloadPlugins(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	pluginDownloadPath := global.Installer.GetAppPluginDownloadPath()
	msg, connectionErr, requestErr := cli.DeleteFiles(pluginDownloadPath)
	if connectionErr != nil {
		ctx.JSON(502, interfaces.Message{Message: connectionErr.Error()})
		return
	}
	ResponseHandler(msg, requestErr, ctx)
}
