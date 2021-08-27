package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"github.com/gin-gonic/gin"
	"github.com/mustafaturan/bus/v3"
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
	userCtx        plugin.UserContext
	config         *Config
	basePath       string
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
	//topic := fmt.Sprintf("%s:%s", "job")
	//eventbus.BUS.RegisterHandler("jobs", BusPluginHandler)
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


var BusPluginHandler = bus.Handler {
	Handle: func(ctx context.Context, e bus.Event) {
		//NewAgent
		data, _ := e.Data.(*model.Point)
		fmt.Println(e.Topic)
		fmt.Println(data)
	},
	Matcher: ".*", // matches all topics
}




// RegisterWebhook implements plugin.Webhooker.
func (c *SystemPlugin) RegisterWebhook(baseURL string, g *gin.RouterGroup) {
	c.basePath = baseURL
	g.GET("/time", func(ctx *gin.Context) {
		ctx.JSON(202, time.Now().Format(time.RFC850))

	})
	g.GET("/time2", func(ctx *gin.Context) {
		ctx.JSON(202, time.Now().Format(time.RFC850))

	})
	g.GET("/time/new", func(ctx *gin.Context) {
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
