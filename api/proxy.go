package api

import (
	"encoding/json"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/client"
	"github.com/gin-gonic/gin"
	"strings"
)

type ProxyDatabase interface {
	GetFN(uuid string) (*model.FlowNetwork, error)
}

type Proxy struct {
	DB ProxyDatabase
}

func (a *Proxy) GetProxy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fnUUID, url := getUrl(ctx)
		fn, err := a.DB.GetFN(fnUUID)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		cli := client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))

		val, err := cli.GetQuery(url)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		response(ctx, val)
	}
}

func (a *Proxy) PostProxy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fnUUID, url := getUrl(ctx)
		fn, err := a.DB.GetFN(fnUUID)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		cli := client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))

		body, err := getMapBody(ctx)
		val, err := cli.PostQuery(url, body)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		response(ctx, val)
	}
}

func (a *Proxy) PutProxy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fnUUID, url := getUrl(ctx)
		fn, err := a.DB.GetFN(fnUUID)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		cli := client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))

		body, err := getMapBody(ctx)
		val, err := cli.PutQuery(url, body)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		response(ctx, val)
	}
}

func (a *Proxy) PatchProxy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fnUUID, url := getUrl(ctx)
		fn, err := a.DB.GetFN(fnUUID)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		cli := client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))

		body, err := getMapBody(ctx)
		val, err := cli.PatchQuery(url, body)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		response(ctx, val)
	}
}

func (a *Proxy) DeleteProxy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fnUUID, url := getUrl(ctx)
		fn, err := a.DB.GetFN(fnUUID)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		cli := client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))

		err = cli.DeleteQuery(url)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		ctx.JSON(204, nil)
	}
}

func getUrl(ctx *gin.Context) (string, string) {
	path := ctx.Request.URL.String()
	subStrings := strings.Split(path, "/")
	fnUUID := subStrings[3]
	url := ""
	for i := 4; i < len(subStrings); i++ {
		url += "/" + subStrings[i]
	}
	return fnUUID, url
}

func response(ctx *gin.Context, val *[]byte) {
	if (*val)[0] == '{' {
		var output map[string]interface{}
		_ = json.Unmarshal(*val, &output)
		ctx.JSON(200, output)
	} else {
		var output []interface{}
		_ = json.Unmarshal(*val, &output)
		ctx.JSON(200, output)
	}
}
