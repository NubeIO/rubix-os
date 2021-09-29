package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
)

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

}
