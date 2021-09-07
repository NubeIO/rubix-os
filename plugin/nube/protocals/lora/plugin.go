package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/handler"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/tty"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
)

// PluginTest is plugin instance
type PluginTest struct {
	config     *Config
	enabled    bool
	msgHandler plugin.MessageHandler
	basePath   string
	UserCtx    plugin.UserContext
	H          handler.Handler
}


func SerialOpenAndRead() {
	bb := new(tty.SerialSetting)
	bb.BaudRate = 9600
	aa := tty.New(bb)
	aa.NewSerialConnection()
	aa.Loop()
}


// Enable implements plugin.Plugin
func (c *PluginTest) Enable() error {
	aaa := c.H.GetNetworks()
	for i, a := range aaa {
		fmt.Println(i, a.Name, "networks")
	}
	//go SerialOpenAndRead()
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
