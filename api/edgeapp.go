package api

import (
	"errors"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/services/systemctl"
	"github.com/NubeIO/flow-framework/src/cli/cligetter"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type EdgeAppDatabase interface {
	ResolveHost(uuid string, name string) (*model.Host, error)
}

type EdgeAppApi struct {
	DB EdgeAppDatabase
}

func (a *EdgeAppApi) EdgeAppUpload(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeClient(host)
	var m *interfaces.AppUpload
	err = ctx.ShouldBindJSON(&m)
	data, err := cli.AppUpload(m)
	ResponseHandler(data, err, ctx)
}

func (a *EdgeAppApi) EdgeAppInstall(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeClient(host)
	var m *systemctl.ServiceFile
	err = ctx.ShouldBindJSON(&m)
	data, err := cli.AppInstall(m)
	ResponseHandler(data, err, ctx)
}

func (a *EdgeAppApi) EdgeAppUninstall(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeClient(host)
	appName := ctx.Query("app_name")
	if appName == "" {
		ResponseHandler(nil, errors.New("app_name can't be empty"), ctx)
		return
	}
	data, err := cli.AppUninstall(appName)
	ResponseHandler(data, err, ctx)
}

func (a *EdgeAppApi) EdgeListAppsStatus(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeClient(host)
	data, err := cli.AppsStatus()
	ResponseHandler(data, err, ctx)
}

func (a *EdgeAppApi) EdgeGetAppStatus(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeClient(host)
	appStatus, connectionErr, requestErr := cli.GetAppStatus(ctx.Param("app_name"))
	if connectionErr != nil {
		ctx.JSON(502, interfaces.Message{Message: connectionErr.Error()})
		return
	}
	ResponseHandler(appStatus, requestErr, ctx)
}
