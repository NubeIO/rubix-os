package api

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/src/cli/cligetter"
	"github.com/gin-gonic/gin"
)

type EdgeConfigDatabase interface {
	ResolveHost(uuid string, name string) (*model.Host, error)
}

type EdgeConfigApi struct {
	DB EdgeConfigDatabase
}

func (a *EdgeConfigApi) EdgeReadConfig(ctx *gin.Context) {
	appName := ctx.Query("app_name")
	configName := ctx.Query("config_name")
	if appName == "" {
		ResponseHandler(nil, errors.New("app_name can not be empty"), ctx)
		return
	}
	if configName == "" {
		configName = "config.yml"
	}
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeClient(host)
	data, err := cli.EdgeReadConfig(appName, configName)
	ResponseHandler(data, err, ctx)
}

func (a *EdgeConfigApi) EdgeWriteConfig(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	var m *interfaces.EdgeConfig
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, ctx)
	}
	cli := cligetter.GetEdgeClient(host)
	data, err := cli.EdgeWriteConfig(m)
	ResponseHandler(data, err, ctx)
}
