package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"net/http"
)

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

//markdown guide
const helpText = `
# LoRa Help Guide

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

type Point struct {
	ObjectType struct {
		Options  interface{} `json:"options"`
		Type     string      `json:"type"`	
		Required bool        `json:"required"`
	} `json:"object_type"`
}

//supportedObjects return all objects that are not bacnet
func supportedObjects() *utils.Array {
	out := utils.NewArray()
	objs := utils.ArrayValues(model.ObjectTypes)
	for _, obj := range objs {
		switch obj {
		case model.ObjectTypes.AnalogInput:
			out.Add(obj)
		case model.ObjectTypes.AnalogOutput:
			out.Add(obj)
		case model.ObjectTypes.AnalogValue:
			out.Add(obj)
		case model.ObjectTypes.BinaryInput:
			out.Add(obj)
		case model.ObjectTypes.BinaryOutput:
			out.Add(obj)
		case model.ObjectTypes.BinaryValue:
			out.Add(obj)
		default:
		}
	}
	return out
}

const (
	help = "/system/help"
	helpHTML = "/system/help/guide"
	pointHelp = "/system/point/help"
)

var Supports = struct {
	Network bool `json:"network"`
	NetworkCRUD bool `json:"networkCRUD"`

}{
	Network: true,
	NetworkCRUD: true,
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.GET(help, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, Supports)
	})
	mux.GET(helpHTML, func(ctx *gin.Context) {
		md := []byte(helpText)
		output := markdown.ToHTML(md, nil, nil)
		ctx.Writer.Write(output)
	})
	mux.GET(pointHelp, func(ctx *gin.Context) {
		var h Point
		h.ObjectType.Options = supportedObjects()
		h.ObjectType.Type = "array"
		h.ObjectType.Required = true
		ctx.JSON(http.StatusOK, h)

	})
	mux.GET("/system/schedule/store/:name", func(ctx *gin.Context) {
		obj, ok := i.store.Get(resolveName(ctx))
		if ok != true {
			ctx.JSON(http.StatusBadRequest, "no schedule exists")
		} else {
			ctx.JSON(http.StatusOK, obj)
		}

	})

}
