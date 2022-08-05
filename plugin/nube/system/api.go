package main

import (
	"github.com/NubeIO/flow-framework/plugin/nube/system/jsonschema"
	"github.com/NubeIO/flow-framework/plugin/nube/system/smodel"
	"github.com/NubeIO/flow-framework/utils/array"
	"net/http"

	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
)

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

const (
	schemaNetwork     = "/schema/network"
	schemaDevice      = "/schema/device"
	schemaPoint       = "/schema/point"
	jsonSchemaNetwork = "/schema/json/network"
	jsonSchemaDevice  = "/schema/json/device"
	jsonSchemaPoint   = "/schema/json/point"
	help              = "/help"
	helpHTML          = "/help/guide"
)

// markdown guide
const helpText = `
# LoRa Help Guide

help

### line 2
*new tab*
<a href="https://stackoverflow.com" target="_blank">New Tab</a>
You will never use anything else than this [website].

- this is some normal texy
- this is ***some*** normal texy
-- aaaaaa
1. First item
2. Second item
3. Third item

| Syntax | Description |
| ----------- | ----------- |
| Header | Title |
| Paragraph | Text |

- [x] Write the press release
- [ ] Update the website
- [ ] Contact the media

this is *some* normal texy`

// supportedObjects return all objects that are not bacnet
func supportedObjects() *array.Array {
	out := array.NewArray()
	out.Add(model.ObjTypeAnalogInput)
	out.Add(model.ObjTypeAnalogOutput)
	out.Add(model.ObjTypeAnalogValue)
	out.Add(model.ObjTypeBinaryInput)
	out.Add(model.ObjTypeBinaryOutput)
	out.Add(model.ObjTypeBinaryValue)
	return out
}

var Supports = struct {
	Network     bool `json:"network"`
	NetworkCRUD bool `json:"networkCRUD"`
}{
	Network:     true,
	NetworkCRUD: true,
}

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	mux.GET(help, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, Supports)
	})
	mux.GET(helpHTML, func(ctx *gin.Context) {
		md := []byte(helpText)
		output := markdown.ToHTML(md, nil, nil)
		ctx.Writer.Write(output)
	})
	mux.GET("/system/schedule/store/:name", func(ctx *gin.Context) {
		obj, ok := inst.store.Get(resolveName(ctx))
		if ok != true {
			ctx.JSON(http.StatusBadRequest, "no schedule exists")
		} else {
			ctx.JSON(http.StatusOK, obj)
		}
	})

	mux.GET(schemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, smodel.GetNetworkSchema())
	})
	mux.GET(schemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, smodel.GetDeviceSchema())
	})
	mux.GET(schemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, smodel.GetPointSchema())
	})
	mux.GET(jsonSchemaNetwork, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, jsonschema.GetNetworkSchema())
	})
	mux.GET(jsonSchemaDevice, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, jsonschema.GetDeviceSchema())
	})
	mux.GET(jsonSchemaPoint, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, jsonschema.GetPointSchema())
	})
}
