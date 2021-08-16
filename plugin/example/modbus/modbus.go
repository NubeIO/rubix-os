package main

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/plugin-api"
	"github.com/gin-gonic/gin"
	"github.com/simonvetter/modbus"
	"log"
	"net/url"
	"time"
)

var (
	client *modbus.ModbusClient
	err    error
)

func startServer() {
	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://192.168.15.202:502",
		Timeout: 1 * time.Second,
	})
}

// GetFlowPluginInfo returns flow plugin info.
func GetFlowPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath: "modbus",
		Name:       "modbus plugin",
	}
}

// EchoPlugin is the flow plugin instance.
type EchoPlugin struct {
	msgHandler     plugin.MessageHandler
	storageHandler plugin.StorageHandler
	userCtx        plugin.UserContext
	config         *Config
	basePath       string
}

// SetStorageHandler implements plugin.Storager
func (c *EchoPlugin) SetStorageHandler(h plugin.StorageHandler) {
	c.storageHandler = h
}

// SetMessageHandler implements plugin.Messenger.
func (c *EchoPlugin) SetMessageHandler(h plugin.MessageHandler) {
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
func (c *EchoPlugin) DefaultConfig() interface{} {
	return &Config{
		MagicString: "hello world",
	}
}

// ValidateAndSetConfig implements plugin.Configurer
func (c *EchoPlugin) ValidateAndSetConfig(config interface{}) error {
	c.config = config.(*Config)
	return nil
}

// Enable enables the plugin.
func (c *EchoPlugin) Enable() error {
	networks, err := c.GetNetworks()
	if err != nil {
		log.Println("GetNetworks GetNetworks ERROR")
	}
	log.Println("GetNetworks GetNetworks GetNetworks", 9999999, networks)
	return nil
}

// GetNetworks disables the plugin.
func (c *EchoPlugin) GetNetworks() ([]*model.Network, error) {
	log.Println("echo plugin GetNetworks")
	fmt.Println(c.userCtx.ID, c.userCtx.Name)
	return nil, nil
}

// GetNetwork disables the plugin.
func (c *EchoPlugin) GetNetwork(id string) error {
	log.Println("echo plugin GetNetworks")
	fmt.Println( c.userCtx.ID, c.userCtx.Name, c.userCtx.Admin)
	return nil
}

// Disable disables the plugin.
func (c *EchoPlugin) Disable() error {
	log.Println("echo plugin disbled")
	fmt.Println( c.userCtx.ID, c.userCtx.Name, c.userCtx.Admin)
	return nil
}

// RegisterWebhook implements plugin.Webhooker.
func (c *EchoPlugin) RegisterWebhook(baseURL string, g *gin.RouterGroup) {
	c.basePath = baseURL
	g.GET("/echo", func(ctx *gin.Context) {
		storage, _ := c.storageHandler.Load()
		//storage
		conf := new(Storage)
		net, err := c.storageHandler.GetNet()
		fmt.Println(net)
		s, _ := json.MarshalIndent(net, "", "\t")
		fmt.Print(string(s))
		fmt.Println(net)
		if err != nil {
			fmt.Println(err)
			//return
		}

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

	g.GET("/echo2", func(ctx *gin.Context) {

		if err != nil {
			// error out if we failed to connect/open the device
			// note: multiple Open() attempts can be made on the same client until
			// the connection succeeds (i.e. err == nil), calling the constructor again
			// is unnecessary.
			// likewise, a client can be opened and closed as many times as needed.
		}

		// read a single 16-bit holding register at address 100
		var reg16 uint16
		reg16, err = client.ReadRegister(0, modbus.HOLDING_REGISTER)
		if err != nil {
			// error out
		} else {
			// use value
			fmt.Println("value: %v", reg16)          // as unsigned integer
			fmt.Println("value: %v", float64(reg16)) // as signed integer
			ctx.JSON(202, reg16)
		}

	})
}

// GetDisplay implements plugin.Displayer.
func (c *EchoPlugin) GetDisplay(location *url.URL) string {
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
	if client == nil {
		startServer()
		err = client.Open()
	}

	return &EchoPlugin{}
}

func main() {
	panic("this should be built as go plugin")
}
