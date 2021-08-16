package main

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"github.com/gin-gonic/gin"
	"time"

	"log"
	"net/url"

)




// GetFlowPluginInfo returns flow plugin info.
func GetFlowPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath: "system",
		Name:       "system time",
	}
}

// SystemPlugin is the flow plugin instance.
type SystemPlugin struct {
	msgHandler     plugin.MessageHandler
	storageHandler plugin.StorageHandler
	userCtx        plugin.UserContext
	config         *Config
	basePath       string
}

// SetStorageHandler implements plugin.Storager
func (c *SystemPlugin) SetStorageHandler(h plugin.StorageHandler) {
	c.storageHandler = h
}

// SetMessageHandler implements plugin.Messenger.
func (c *SystemPlugin) SetMessageHandler(h plugin.MessageHandler) {
	c.msgHandler = h
}

// Storage defines the plugin storage scheme
type Storage struct {
	CalledTimes int `json:"called_times"`
}

// Config defines the plugin config scheme
type Config struct {
	MagicString string `yaml:"magic_string"`
}

// DefaultConfig implements plugin.Configurer
func (c *SystemPlugin) DefaultConfig() interface{} {
	return &Config{
		MagicString: "hello world",
	}
}

// ValidateAndSetConfig implements plugin.Configurer
func (c *SystemPlugin) ValidateAndSetConfig(config interface{}) error {
	c.config = config.(*Config)
	return nil
}

// Enable enables the plugin.
func (c *SystemPlugin) Enable() error {
	log.Println("plugin Enable")
	return nil
}


// Disable disables the plugin.
func (c *SystemPlugin) Disable() error {
	log.Println("plugin Disable")
	return nil
}


// GetNetworks disables the plugin.
func (c *SystemPlugin) GetNetworks() ([]*model.Network, error) {
	log.Println("GetNetworks")
	fmt.Println(c.userCtx.ID, c.userCtx.Name)
	return nil, nil
}

// GetNetwork disables the plugin.
func (c *SystemPlugin) GetNetwork(id string) error {
	log.Println("plugin GetNetworks")
	return nil
}


// RegisterWebhook implements plugin.Webhooker.
func (c *SystemPlugin) RegisterWebhook(baseURL string, g *gin.RouterGroup) {
	c.basePath = baseURL
	g.GET("/message", func(ctx *gin.Context) {
		storage, _ := c.storageHandler.Load()
		//storage
		conf := new(Storage)
		json.Unmarshal(storage, conf)
		conf.CalledTimes++
		newStorage, _ := json.Marshal(conf)
		c.storageHandler.Save(newStorage)
		c.msgHandler.SendMessage(plugin.Message{
			Title:    "Hello received",
			Message:  fmt.Sprintf("echo server received a hello message %d times", conf.CalledTimes),
			Priority: 2,
			Extras: map[string]interface{}{
				"plugin::name": "echo",
			},
		})
		ctx.Writer.WriteString(fmt.Sprintf("Magic string is: %s\r\nEcho server running at %secho", c.config.MagicString, c.basePath))
	})

	g.GET("/time", func(ctx *gin.Context) {
		ctx.JSON(202, time.Now().Format(time.RFC850))

	})

}

// GetDisplay implements plugin.Displayer.
func (c *SystemPlugin) GetDisplay(location *url.URL) string {
	loc := &url.URL{
		Path: c.basePath,
	}
	if location != nil {
		loc.Scheme = location.Scheme
		loc.Host = location.Host
	}
	loc = loc.ResolveReference(&url.URL{
		Path: "echo",
	})
	return "Echo plugin running at: " + loc.String()
}

// NewFlowPluginInstance creates a plugin instance for a user context.
func NewFlowPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &SystemPlugin{}
}

func main() {
	panic("this should be built as go plugin")
}
