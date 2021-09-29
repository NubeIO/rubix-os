package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"github.com/NubeDev/flow-framework/utils"
	"net/url"
)

//supportedObjects return all objects that are not bacnet
func supportedObjects() *utils.Array {
	out := utils.NewArray()
	objs := utils.ArrayValues(model.ObjectTypes)
	for _, obj := range objs {
		switch obj {
		case model.ObjectTypes.AnalogInput:
			out.Add(obj)
		case model.ObjectTypes.AnalogOutput:
			out.Add(obj)
		case model.ObjectTypes.AnalogValue:
			out.Add(obj)
		case model.ObjectTypes.BinaryInput:
			out.Add(obj)
		case model.ObjectTypes.BinaryOutput:
			out.Add(obj)
		case model.ObjectTypes.BinaryValue:
			out.Add(obj)
		default:
		}
	}
	return out
}

// GetDisplay implements public.Displayer
func (i *Instance) GetDisplay(baseURL *url.URL) plugin.Response {
	loc := &url.URL{
		Path: i.basePath,
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "restart",
	})
	fmt.Println(loc) //can show the ui the custom endpoints

	baseURL.Path = i.basePath
	m := plugin.Help{
		Name:               name,
		PluginType:         pluginType,
		AllowConfigWrite:   allowConfigWrite,
		IsNetwork:          isNetwork,
		MaxAllowedNetworks: maxAllowedNetworks,
		NetworkType:        networkType,
		TransportType:      transportType,
	}
	messageURL := plugin.Response{
		Details: m,
	}
	return messageURL
}
