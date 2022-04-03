package main

import (
	"github.com/NubeIO/flow-framework/plugin/pluginapi"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// PluginTest is plugin instance
type PluginTest struct {
	config     *Config
	enabled    bool
	msgHandler pluginapi.MessageHandler
	basePath   string
}

func (c *PluginTest) GetNetworks() ([]*model.Network, error) {
	panic("implement me")
}

func (c *PluginTest) GetNetwork(id string) error {
	panic("implement me")
}

// Enable implements plugin.Plugin
func (c *PluginTest) Enable() error {
	c.enabled = true
	return nil
}

// Disable implements plugin.Disable
func (c *PluginTest) Disable() error {
	c.enabled = false
	return nil
}

// SetMessageHandler implements plugin.Messenger
func (c *PluginTest) SetMessageHandler(h pluginapi.MessageHandler) {
	c.msgHandler = h
}
