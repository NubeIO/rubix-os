package main

import (
	"github.com/NubeIO/flow-framework/plugin/plugin-api"
	"net/url"
)

// GetDisplay implements public.Displayer
func (c *PluginTest) GetDisplay(baseURL *url.URL) plugin.Response {
	baseURL.Path = c.basePath
	messageURL := plugin.Response{
		Details: plugin.Help{
			Name: "System",
		},
	}
	return messageURL
}
