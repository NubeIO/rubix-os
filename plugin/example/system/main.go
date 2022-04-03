package main

import (
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
)

const path = "system"
const name = "system"
const description = "system plugin"
const author = "ap"
const webSite = "https://www.github.com/NubeIO"

// GetFlowPluginInfo returns plugin info.
func GetFlowPluginInfo() pluginapi.Info {
	return pluginapi.Info{
		ModulePath:  path,
		Name:        name,
		Description: description,
		Author:      author,
		Website:     webSite,
	}
}

// NewFlowPluginInstance creates a plugin instance for a user context.
func NewFlowPluginInstance() pluginapi.Plugin {
	p := &PluginTest{}
	return p
}

//main will not let main run
func main() {
	panic("this should be built as plugin")
}
