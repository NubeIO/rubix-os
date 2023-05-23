package main

import (
	"net/url"

	"github.com/NubeIO/rubix-os/plugin/pluginapi"
)

func (inst *Instance) GetDisplay(baseURL *url.URL) pluginapi.Response {
	loc := &url.URL{
		Path: inst.basePath,
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "restart",
	})

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
