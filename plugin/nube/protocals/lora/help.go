package main

import (
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"net/url"
)

// GetDisplay implements public.Displayer
func (c *Instance) GetDisplay(baseURL *url.URL) plugin.Response {
	baseURL.Path = c.basePath
	m := plugin.Help{
		Name: name,
		PluginType: pluginType,
		AllowConfigWrite: allowConfigWrite,
		IsNetwork: isNetwork,
		MaxAllowedNetworks: maxAllowedNetworks,
		NetworkType: networkType,
		TransportType: transportType,
	}
	messageURL := plugin.Response {
		Details:       m,
	}
	return messageURL
}
