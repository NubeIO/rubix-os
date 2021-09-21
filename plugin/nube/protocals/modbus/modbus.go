package main

import (
	"github.com/NubeDev/flow-framework/utils"
	"github.com/simonvetter/modbus"
	"time"
)

type Client struct {
	Host    string        `json:"ip"`
	Port    string        `json:"port"`
	Timeout time.Duration `json:"device_timeout_in_ms"`
}

var restMB *modbus.ModbusClient
var connected bool

func setClient(client Client) error {
	var cli utils.URLParts
	cli.Transport = "tcp"
	cli.Host = client.Host
	cli.Port = client.Port
	url := utils.JoinURL(cli)

	if client.Timeout < 10 {
		client.Timeout = 500
	}
	c, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     url,
		Timeout: client.Timeout * time.Millisecond,
	})
	if err != nil {
		connected = false
		return err
	}
	connected = true
	err = c.Open()
	restMB = c
	if err != nil {
		connected = false
		return err
	}
	return nil
}

func getClient() *modbus.ModbusClient {
	return restMB
}

func isConnected() bool {
	return connected
}
