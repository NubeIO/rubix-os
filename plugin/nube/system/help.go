package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"net/url"
)

// GetDisplay implements public.Displayer
func (i *Instance) GetDisplay(baseURL *url.URL) pluginapi.Response {
	loc := &url.URL{
		Path: i.basePath,
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "restart",
	})
	fmt.Println(loc) //can show the ui the custom endpoints

	baseURL.Path = i.basePath
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
