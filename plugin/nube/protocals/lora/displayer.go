package main

import (
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"net/url"
)





// GetDisplay implements public.Displayer
func (c *PluginTest) GetDisplay(baseURL *url.URL) plugin.Response {
	baseURL.Path = c.basePath
	messageURL := plugin.Response {
		StatusCode: 1,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       "Hello World",
	}
	return messageURL
}
