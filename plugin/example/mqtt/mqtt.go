package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/flow-framework/plugin/plugin-api"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

var (
	ErrInvalidAddress = errors.New("invalid broker address")
)

// GetFlowPluginInfo returns Flow plugin info
func GetFlowPluginInfo() plugin.Info {
	return plugin.Info{
		Name:       "MQTT",
		ModulePath: "Flow-mqttClient",
		Author:     "ap",
		Website:    "nube-io.com",
	}
}

type Server struct {
	Address   string
	Username  string
	Password  string
	Subscribe []string
}

type Config struct {
	Servers []Server
}

// Plugin is plugin instance
type Plugin struct {
	msgHandler plugin.MessageHandler
	config     *Config
	clients    []mqtt.Client
	enabled    bool
}

func (p *Plugin) GetNetworks() ([]*model.Network, error) {
	return nil, nil
}

func (p *Plugin) GetNetwork(id string) error {
	return nil
}

// SetMessageHandler implements plugin.Messenger
// Invoked during initialization
func (p *Plugin) SetMessageHandler(h plugin.MessageHandler) {
	p.msgHandler = h
}

// Enable adds users to the context map which maps to a Plugin.
func (p *Plugin) Enable() error {
	p.enabled = true
	return nil
}

// Disable removes users from the context map.
func (p *Plugin) Disable() error {
	p.enabled = false
	p.disconnectClients()
	return nil
}

// DefaultConfig implements plugin.Configurer
func (p *Plugin) DefaultConfig() interface{} {
	return &Config{
		Servers: []Server{
			{Address: "127.0.0.1:1883", Subscribe: []string{"*"}},
		},
	}
}

// ValidateAndSetConfig will be called every time the plugin is initialized or the configuration has been changed by the user.
func (p *Plugin) ValidateAndSetConfig(c interface{}) error {
	config := c.(*Config)

	// If listeners are configured, shut them down and start fresh
	if p.clients != nil {
		for _, client := range p.clients {
			if client == nil || !client.IsConnected() {
				continue
			}

			go client.Disconnect(500)
		}
	}
	p.clients = make([]mqtt.Client, len(config.Servers))
	for _, server := range config.Servers {
		if server.Address == "" {
			return ErrInvalidAddress
		}
	}

	p.config = config
	// If enabled already and config was updated, reconnect clients
	if p.enabled {
		return p.connectClients()
	}

	return nil
}

func (p *Plugin) disconnectClients() {
	if p.clients == nil {
		return
	}

	for _, client := range p.clients {
		if client == nil || !client.IsConnected() {
			continue
		}

		go client.Disconnect(500)
	}
}

func (p *Plugin) connectClients() error {
	p.disconnectClients()

	p.clients = make([]mqtt.Client, len(p.config.Servers))

	for i, server := range p.config.Servers {
		client, err := p.newClient(server)

		if err != nil {
			return err
		}

		p.clients[i] = client
	}

	return nil
}

// RegisterWebhook implements plugin.Webhooker.
func (p *Plugin) RegisterWebhook(baseURL string, g *gin.RouterGroup) {
	g.POST("/mqttClient", func(ctx *gin.Context) {
		for _, a := range p.clients {
			log.Println("try a message")
			if a.IsConnected() {
				log.Println("send a message")
				msg := fmt.Sprintf("hello from MQTT %s time", time.Now().Format(time.RFC850))
				a.Publish("topic", 1, false, msg)
				err := p.msgHandler.SendMessage(plugin.Message{
					Title:    "mqttClient-message",
					Message:  fmt.Sprintf("hello from rest %s time", time.Now().Format(time.RFC850)),
					Priority: 2,
					Extras: map[string]interface{}{
						"plugin::name": "echo",
					},
				})
				if err != nil {
					ctx.JSON(404, "error on send message")
				}
				rmsg := fmt.Sprintf("hello from REST %s time", time.Now().Format(time.RFC850))
				ctx.JSON(202, rmsg)
			} else {
				ctx.JSON(404, "broker not connected")
			}
		}
	})
}

// handleMessage handles mqtt messages from the client by returning a MessageHandler
func (p *Plugin) handleMessage(client mqtt.Client, message mqtt.Message) {
	payload := message.Payload()

	var outgoingMessage plugin.Message

	if payload[0] == '{' {
		if err := json.Unmarshal(payload, &outgoingMessage); err != nil {
			return
		}
	} else {
		outgoingMessage.Message = string(payload)
	}

	p.msgHandler.SendMessage(outgoingMessage)
}

// newClient creates a new client from the serverConfig
func (p *Plugin) newClient(serverConfig Server) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(serverConfig.Address)
	opts.SetClientID("Flow")

	if serverConfig.Username != "" {
		opts.SetUsername(serverConfig.Username)
	}

	if serverConfig.Password != "" {
		opts.SetPassword(serverConfig.Password)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	for _, topic := range serverConfig.Subscribe {
		client.Subscribe(topic, 0, p.handleMessage)
	}

	return client, nil
}

// NewFlowPluginInstance creates a plugin instance for a user context.
func NewFlowPluginInstance() plugin.Plugin {
	//return &Plugin{}
	return &Plugin{
		clients: make([]mqtt.Client, 0),
	}
}

func main() {
	panic("Program must be compiled as a Go plugin")
}
