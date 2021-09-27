package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/simonvetter/modbus"
	"time"
)

type Client struct {
	Host       string        `json:"ip"`
	Port       string        `json:"port"`
	SerialPort string        `json:"serial_port"`
	BaudRate   uint          `json:"baud_rate"`
	StopBits   uint          `json:"stop_bits"`
	Parity     string        `json:"parity"` //none, odd, even
	DataBits   uint          `json:"data_bits"`
	Timeout    time.Duration `json:"device_timeout_in_ms"`
}

var restMB *modbus.ModbusClient
var connected bool

func setClient(client Client) error {
	var cli utils.URLParts
	cli.Transport = "tcp"
	cli.Host = client.Host
	cli.Port = client.Port
	url, err := utils.JoinURL(cli)
	if err != nil {
		connected = false
		return err
	}
	if client.Timeout < 10 {
		client.Timeout = 500
	}
	//TODO add in a check if client with same details exists
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

func setParity(in string) uint {
	if in == model.SerialParity.None {
		return modbus.PARITY_NONE
	} else if in == model.SerialParity.Odd {
		return modbus.PARITY_ODD
	} else if in == model.SerialParity.Even {
		return modbus.PARITY_EVEN
	} else {
		return modbus.PARITY_NONE
	}

}

func setSerial(port string) string {
	p := fmt.Sprintf("rtu:///%s", port)
	return p
}

func setClientSerial(client Client) error {
	parity := setParity(client.Parity)
	serialPort := setSerial(client.SerialPort)
	if client.Timeout < 10 {
		client.Timeout = 500
	}
	//TODO add in a check if client with same details exists
	c, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      serialPort,
		Speed:    client.BaudRate, // default
		DataBits: client.DataBits, // default, optional
		Parity:   parity,          // default, optional
		StopBits: 2,               // default if no parity, optional
		Timeout:  client.Timeout * time.Millisecond,
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
