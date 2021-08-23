package main

import (
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
)

const gitHubURL = "https://www.github.com/NubeDev/flow-framework"


// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:  "github.com/NubeDev/flow-framework/message",
		Name:        "Test",
		Description: "A plugin",
		Author:      "ap",
		Website:     gitHubURL,
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