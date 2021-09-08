package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"net/url"
)


// GetDisplay implements public.Displayer
func (c *Instance) GetDisplay(baseURL *url.URL) plugin.Response {
	loc := &url.URL{
		Path: c.basePath,
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "restart",
	})
	fmt.Println(loc) //can show the ui the custom endpoints

	baseURL.Path = c.basePath
	m := plugin.Help{
		Name:               name,
		PluginType:         pluginType,
		AllowConfigWrite:   allowConfigWrite,
		IsNetwork:          isNetwork,
		MaxAllowedNetworks: maxAllowedNetworks,
		NetworkType:        networkType,
		TransportType:      transportType,
	}
	messageURL := plugin.Response {
		Details:       m,
	}
	return messageURL
}
