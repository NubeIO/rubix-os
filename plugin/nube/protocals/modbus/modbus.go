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
	Parity     string        `json:"parity"` //none, odd, even DEFAULT IS PARITY_NONE
	DataBits   uint          `json:"data_bits"`
	Timeout    time.Duration `json:"device_timeout_in_ms"`
}

var restMB *modbus.ModbusClient
var connected bool

func (i *Instance) setClient(client Client, networkUUID string, cacheClient, isSerial bool) error {
	var c *modbus.ModbusClient
	if isSerial {
		parity := setParity(client.Parity)
		serialPort := setSerial(client.SerialPort)
		if client.Timeout < 10 {
			client.Timeout = 500
		}
		//TODO add in a check if client with same details exists
		c, err = modbus.NewClient(&modbus.ClientConfiguration{
			URL:      serialPort,
			Speed:    client.BaudRate,
			DataBits: client.DataBits,
			Parity:   parity,
			StopBits: client.StopBits,
			Timeout:  client.Timeout * time.Millisecond,
		})
	} else {
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
		c, err = modbus.NewClient(&modbus.ClientConfiguration{
			URL:     url,
			Timeout: client.Timeout * time.Millisecond,
		})
		if err != nil {
			connected = false
			return err
		}
	}
	var getC interface{}
	if cacheClient { //store modbus client in cache to reuse the instance
		getC, _ = i.store.Get(networkUUID)
		if getC == nil {
			i.store.Set(networkUUID, c, -1)
		} else {
			c = getC.(*modbus.ModbusClient)
		}
	}
	err = c.Open()
	connected = true
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
