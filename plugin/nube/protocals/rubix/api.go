package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/rubix/rubixapi"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/rubix/rubixmodel"
	"github.com/gin-gonic/gin"
	"net/http"
)

func bodyToken(ctx *gin.Context) (dto rubixmodel.TokenBody, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyAppControl(ctx *gin.Context) (dto rubixmodel.AppControl, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyAppsDownload(ctx *gin.Context) (dto rubixmodel.AppsDownload, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyJson(ctx *gin.Context) (dto interface{}, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

//user
const (
	rubix     = "rubix"
	users     = "users"
	usersList = "rubix/users"
)

//app
const (
	apps         = "apps"
	appsList     = "rubix/apps"
	appsDownload = "rubix/apps/download"
)

var endPoints = struct {
	Users              string
	AppsControl        string
	AppsInstalled      string
	AppsLatestVersions string
}{
	Users:              fmt.Sprintf("/%s/:name/%s", rubix, users),
	AppsControl:        fmt.Sprintf("/%s/:name/%s/%s", rubix, apps, "control"),
	AppsInstalled:      fmt.Sprintf("/%s/:name/%s", rubix, apps),
	AppsLatestVersions: fmt.Sprintf("/%s/:name/%s/%s", rubix, apps, "latest_versions"),
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath

	mux.GET(endPoints.Users, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		fmt.Println(endPoints.AppsControl)
		req, err := i.getIntegration("", _name)
		r, err := cli.GetUsers(req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, r)
		}
	})

	mux.POST(endPoints.AppsControl, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		body, _ := bodyAppControl(ctx)
		req.Body = body
		r, err := cli.AppControl(req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, r)
		}
	})

	mux.GET(endPoints.AppsInstalled, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.AppsInstalled(req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, r)
		}
	})

	mux.GET(endPoints.AppsLatestVersions, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.AppsLatestVersions(req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			ctx.JSON(http.StatusOK, r)
		}
	})
}
