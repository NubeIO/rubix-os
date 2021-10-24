package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"net/http"
)

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

//markdown guide
const help = `
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

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	//restart plugin
	mux.GET("/system/help", func(ctx *gin.Context) {
		md := []byte(help)
		output := markdown.ToHTML(md, nil, nil)
		ctx.Writer.Write(output)
		//ctx.Writer.Write(output)
		//ctx.Writer.WriteString(fmt.Sprintf("Magic string is: %s\r\nEcho server running at %secho", "22", i.basePath))

	})
	/*
			get the schedule by its name
		    "weekly": {
		            "cf50cd39-e1cf-4d7e-aa70-2dc7220780f1": {
		                "name": "Branch",
	*/
	mux.GET("/system/schedule/store/:name", func(ctx *gin.Context) {
		obj, ok := i.store.Get(resolveName(ctx))
		if ok != true {
			ctx.JSON(http.StatusBadRequest, "no schedule exists")
		} else {
			ctx.JSON(http.StatusOK, obj)
		}

	})

}
