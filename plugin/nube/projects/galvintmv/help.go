package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"net/url"
)

func (inst *Instance) GetDisplay(baseURL *url.URL) pluginapi.Response {
	loc := &url.URL{
		Path: inst.basePath,
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "restart",
	})
	fmt.Println(loc) // can show the ui the custom endpoints

	baseURL.Path = inst.basePath
	m := pluginapi.Help{
		Name:               name,
		PluginType:         pluginType,
		AllowConfigWrite:   allowConfigWrite,
		IsNetwork:          isNetwork,
		MaxAllowedNetworks: maxAllowedNetworks,
		NetworkType:        networkType,
		TransportType:      transportType,
	}
	messageURL := pluginapi.Response{
		Details: m,
	}
	return messageURL
}
