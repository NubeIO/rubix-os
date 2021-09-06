package main

import (
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
)

const path = "system"
const name = "system"
const description = "system plugin"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:  path,
		Name:        name,
		Description: description,
		Author:      author,
		Website:     webSite,
	}
}

// NewFlowPluginInstance creates a plugin instance for a user context.
func NewFlowPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	usersList.AddUser(ctx)
	p := &PluginTest{
		UserCtx: ctx,
	}
	return p
}

//main will not let main run
func main() {
	panic("this should be built as plugin")
}