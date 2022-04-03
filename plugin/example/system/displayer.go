package main

import (
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"net/url"
)

// GetDisplay implements public.Displayer
func (c *PluginTest) GetDisplay(baseURL *url.URL) pluginapi.Response {
	baseURL.Path = c.basePath
	messageURL := pluginapi.Response{
		Details: pluginapi.Help{
			Name: "System",
		},
	}
	return messageURL
}
