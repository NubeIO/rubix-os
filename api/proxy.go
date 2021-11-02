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
	GetFNC(uuid string) (*model.FlowNetworkClone, error)
}

type Proxy struct {
	DB ProxyDatabase
}

func (a *Proxy) GetProxy(isFN bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid, url := getUrl(ctx)
		cli, err := a.getFlowClient(isFN, uuid)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}

		val, err := cli.GetQuery(url)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		response(ctx, val)
	}
}

func (a *Proxy) PostProxy(isFN bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid, url := getUrl(ctx)
		cli, err := a.getFlowClient(isFN, uuid)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}

		body, err := getMapBody(ctx)
		val, err := cli.PostQuery(url, body)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		response(ctx, val)
	}
}

func (a *Proxy) PutProxy(isFN bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid, url := getUrl(ctx)
		cli, err := a.getFlowClient(isFN, uuid)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}

		body, err := getMapBody(ctx)
		val, err := cli.PutQuery(url, body)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		response(ctx, val)
	}
}

func (a *Proxy) PatchProxy(isFN bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid, url := getUrl(ctx)
		cli, err := a.getFlowClient(isFN, uuid)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}

		body, err := getMapBody(ctx)
		val, err := cli.PatchQuery(url, body)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		response(ctx, val)
	}
}

func (a *Proxy) DeleteProxy(isFN bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuid, url := getUrl(ctx)
		cli, err := a.getFlowClient(isFN, uuid)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}

		err = cli.DeleteQuery(url)
		if err != nil {
			ctx.AbortWithError(404, err)
			return
		}
		ctx.JSON(204, nil)
	}
}

func (a *Proxy) getFlowClient(isFN bool, uuid string) (*client.FlowClient, error) {
	var cli *client.FlowClient
	if isFN {
		fn, err := a.DB.GetFN(uuid)
		if err != nil {
			return nil, err
		}
		cli = client.NewFlowClientCli(fn.FlowIP, fn.FlowPort, fn.FlowToken, fn.IsMasterSlave, fn.GlobalUUID, model.IsFNCreator(fn))
	} else {
		fnc, err := a.DB.GetFNC(uuid)
		if err != nil {
			return nil, err
		}
		cli = client.NewFlowClientCli(fnc.FlowIP, fnc.FlowPort, fnc.FlowToken, fnc.IsMasterSlave, fnc.GlobalUUID, model.IsFNCreator(fnc))
	}
	return cli, nil
}

func getUrl(ctx *gin.Context) (string, string) {
	path := ctx.Request.URL.String()
	subStrings := strings.Split(path, "/")
	uuid := subStrings[3]
	url := ""
	for i := 4; i < len(subStrings); i++ {
		url += "/" + subStrings[i]
	}
	return uuid, url
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
