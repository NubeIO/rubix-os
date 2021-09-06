package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/tty"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
)

// PluginTest is plugin instance
type PluginTest struct {
	config     *Config
	enabled    bool
	msgHandler plugin.MessageHandler
	basePath   string

	UserCtx plugin.UserContext
}

func (c *PluginTest) GetNetworks() ([]*model.Network, error) {
	panic("implement me")
}

func (c *PluginTest) GetNetwork(id string) error {
	panic("implement me")
}

func SerialOpenAndRead() {
	bb := new(tty.SerialSetting)
	bb.BaudRate = 38400
	aa := tty.New(bb)
	aa.NewSerialConnection()
	aa.Loop()
}


// Enable implements plugin.Plugin
func (c *PluginTest) Enable() error {
	go SerialOpenAndRead()
	c.enabled = true
	return nil
}

// Disable implements plugin.Disable
func (c *PluginTest) Disable() error {
	c.enabled = false
	return nil
}

// SetMessageHandler implements plugin.Messenger
func (c *PluginTest) SetMessageHandler(h plugin.MessageHandler) {
	c.msgHandler = h
}
