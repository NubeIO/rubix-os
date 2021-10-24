package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/rubix/rubixapi"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/rubix/rubixmodel"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

func bodyAppControl(ctx *gin.Context) (dto rubixmodel.AppControl, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyAppsDownload(ctx *gin.Context) (dto rubixmodel.AppsDownload, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyWiresPlat(ctx *gin.Context) (dto rubixmodel.WiresPlat, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyGlobalUUID(ctx *gin.Context) (dto rubixmodel.GeneralResponse, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

const (
	rubix    = "rubix"
	users    = "users"
	apps     = "apps"
	discover = "discover"
	slaves   = "slaves"
	wires    = "wires"
)

func httpRes(obj interface{}, err error, ctx *gin.Context) {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	} else {
		ctx.JSON(http.StatusOK, obj)
	}
}

func urlPath(u string) (clean string) {
	_url := fmt.Sprintf("http://%s", u)
	p, err := url.Parse(_url)
	if err != nil {
		return ""
	}
	parts := strings.SplitAfter(p.String(), "any")
	if len(parts) >= 1 {
		return parts[1]
	} else {
		return ""
	}
}

var endPoints = struct {
	Users                 string
	AppsControl           string
	AppsInstalled         string
	AppsLatestVersions    string
	AppsInstall           string
	AppsDownloadState     string
	AppsReleases          string
	DiscoverRemoteDevices string
	Slaves                string
	SlavesDelete          string
	WiresPlat             string
	Proxy                 string
}{
	Users:                 fmt.Sprintf("/%s/:name/%s", rubix, users),
	AppsControl:           fmt.Sprintf("/%s/:name/%s/%s", rubix, apps, "control"),
	AppsInstalled:         fmt.Sprintf("/%s/:name/%s", rubix, apps),
	AppsLatestVersions:    fmt.Sprintf("/%s/:name/%s/%s", rubix, apps, "latest_versions"),
	AppsInstall:           fmt.Sprintf("/%s/:name/%s/%s", rubix, apps, "install"),
	AppsDownloadState:     fmt.Sprintf("/%s/:name/%s/%s", rubix, apps, "download_state"),
	AppsReleases:          fmt.Sprintf("/%s/:name/%s/%s", rubix, apps, "releases"),
	DiscoverRemoteDevices: fmt.Sprintf("/%s/:name/%s/%s", rubix, discover, "remote_devices"),
	Slaves:                fmt.Sprintf("/%s/:name/%s", rubix, slaves),
	SlavesDelete:          fmt.Sprintf("/%s/:name/%s/:uuid", rubix, slaves),
	WiresPlat:             fmt.Sprintf("/%s/:name/%s/%s", rubix, wires, "plat"),
	Proxy:                 fmt.Sprintf("/%s/:name/proxy/*proxy", rubix),
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath

	mux.GET(endPoints.Users, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.GetUsers(req)
		httpRes(r, err, ctx)
	})
	/*
		APPS
	*/
	mux.POST(endPoints.AppsControl, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		body, _ := bodyAppControl(ctx)
		req.Body = body
		r, err := cli.AppControl(req)
		httpRes(r, err, ctx)
	})

	mux.GET(endPoints.AppsInstalled, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.AppsInstalled(req)
		httpRes(r, err, ctx)
	})

	mux.GET(endPoints.AppsLatestVersions, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.AppsLatestVersions(req)
		httpRes(r, err, ctx)
	})

	mux.POST(endPoints.AppsInstall, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		body, _ := bodyAppsDownload(ctx)
		req.Body = body
		r, err := cli.AppsInstall(req)
		httpRes(r, err, ctx)
	})

	mux.GET(endPoints.AppsDownloadState, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.AppsDownloadState(req)
		httpRes(r, err, ctx)
	})

	mux.DELETE(endPoints.AppsDownloadState, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.AppsDeleteDownloadState(req)
		httpRes(r, err, ctx)
	})
	/*
		SLAVES, discover get, add delete
	*/
	mux.GET(endPoints.DiscoverRemoteDevices, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.SlaveDevices(req, true)
		httpRes(r, err, ctx)
	})
	mux.GET(endPoints.Slaves, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.SlaveDevices(req, false)
		httpRes(r, err, ctx)
	})
	mux.POST(endPoints.Slaves, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		body, _ := bodyGlobalUUID(ctx)
		req.Body = body
		r, err := cli.SlaveDevicesAddDelete(req, false, body.GlobalUUID)
		httpRes(r, err, ctx)
	})
	mux.DELETE(endPoints.Slaves, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		body, _ := bodyGlobalUUID(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.SlaveDevicesAddDelete(req, true, body.GlobalUUID)
		httpRes(r, err, ctx)
	})
	/*
		WIRES_PLAT, get, add edit
	*/
	mux.GET(endPoints.WiresPlat, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		r, err := cli.WiresPlat(req, false)
		httpRes(r, err, ctx)
	})
	mux.PUT(endPoints.WiresPlat, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		body, _ := bodyWiresPlat(ctx)
		req, err := i.getIntegration("", _name)
		req.Body = body
		r, err := cli.WiresPlat(req, true)
		httpRes(r, err, ctx)
	})
	/*
		Proxy
		Will get the incoming url path and body and forward on the request
	*/
	mux.GET(endPoints.Proxy, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			req.URL = urlPath(ctx.Request.URL.String())
			if req.URL == "" || req.URL == "/" {
				ctx.JSON(http.StatusBadRequest, "invalid request")
			} else {
				r, err := cli.AnyRequest(req)
				req.Method = rubixapi.GET
				httpRes(r, err, ctx)
			}
		}

	})
	mux.DELETE(endPoints.Proxy, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		req, err := i.getIntegration("", _name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			req.URL = urlPath(ctx.Request.URL.String())
			if req.URL == "" || req.URL == "/" {
				ctx.JSON(http.StatusBadRequest, "invalid request")
			} else {
				r, err := cli.AnyRequest(req)
				req.Method = rubixapi.DELETE
				httpRes(r, err, ctx)
			}
		}

	})
	mux.POST(endPoints.Proxy, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		var getBody interface{} //get the body and put it into an interface
		err := ctx.ShouldBindJSON(&getBody)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		}
		req, err := i.getIntegration("", _name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			req.URL = urlPath(ctx.Request.URL.String())
			if req.URL == "" || req.URL == "/" {
				ctx.JSON(http.StatusBadRequest, "invalid request")
			} else {
				req.URL = urlPath(ctx.Request.URL.String()) //pass on the path
				req.Body = getBody                          //pass on the body
				req.Method = rubixapi.POST
				r, err := cli.AnyRequestWithBody(req)
				httpRes(r, err, ctx)
			}
		}
	})
	mux.PATCH(endPoints.Proxy, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		var getBody interface{} //get the body and put it into an interface
		err := ctx.ShouldBindJSON(&getBody)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		}
		req, err := i.getIntegration("", _name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			req.URL = urlPath(ctx.Request.URL.String())
			if req.URL == "" || req.URL == "/" {
				ctx.JSON(http.StatusBadRequest, "invalid request")
			} else {
				req.URL = urlPath(ctx.Request.URL.String()) //pass on the path
				req.Body = getBody                          //pass on the body
				req.Method = rubixapi.PATCH
				r, err := cli.AnyRequestWithBody(req)
				httpRes(r, err, ctx)
			}
		}
	})
	mux.PUT(endPoints.Proxy, func(ctx *gin.Context) {
		cli := rubixapi.New()
		_name := resolveName(ctx)
		var getBody interface{} //get the body and put it into an interface
		err := ctx.ShouldBindJSON(&getBody)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		}
		req, err := i.getIntegration("", _name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			req.URL = urlPath(ctx.Request.URL.String())
			if req.URL == "" || req.URL == "/" {
				ctx.JSON(http.StatusBadRequest, "invalid request")
			} else {
				req.URL = urlPath(ctx.Request.URL.String()) //pass on the path
				req.Body = getBody                          //pass on the body
				req.Method = rubixapi.PUT
				r, err := cli.AnyRequestWithBody(req)
				httpRes(r, err, ctx)
			}
		}
	})
}
