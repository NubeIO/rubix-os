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
	messageURL := pluginapi.Response{
		Details: pluginapi.Help{
			Name:             pluginName,
			PluginType:       pluginType,
			AllowConfigWrite: allowConfigWrite,
			IsNetwork:        hasNetwork,
			NetworkType:      networkType,
			TransportType:    transportType,
		},
	}
	return messageURL
}
