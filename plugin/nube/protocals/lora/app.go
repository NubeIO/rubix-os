package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/tty"
	log "github.com/sirupsen/logrus"
)

/*
user adds a network
user adds a device
- create device and send plugin the uuid
- ask the plugin do you want to add pre-made points for example
- add points
*/

func SerialOpenAndRead() error {
	s := new(tty.SerialSetting)
	s.SerialPort = "/dev/ttyACM0"
	s.BaudRate = 38400
	sc := tty.New(s)
	err := sc.NewSerialConnection()
	if err != nil {
		return err
	}
	sc.Loop()
	return nil
}

// SerialOpen open serial port
func (c *Instance) SerialOpen() error {
	go func() error {
		err := SerialOpenAndRead()
		if err != nil {
			return err
		}
		return nil
	}()
	log.Info("LORA: open serial port")
	return nil
}

// SerialClose close serial port
func (c *Instance) SerialClose() error {
	err := tty.Disconnect()
	if err != nil {
		return err
	}
	return nil
}

var THLM = []string{"rssi", "voltage", "temperature", "humidity", "light", "motion"}

// addPoints close serial port
func (c *Instance) addPoints(deviceBody *model.Device) (*model.Point, error) {
	p := new(model.Point)
	p.DeviceUUID = deviceBody.UUID
	code := deviceBody.AddressCode
	if code == string(decoder.THLM) {
		for _, e := range THLM {
			p.PointType = e
			err := c.addPoint(p)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil

}

// addPoints close serial port
func (c *Instance) addPoint(body *model.Point) error {
	_, err := c.db.CreatePoint(body)
	if err != nil {
		return err
	}
	return nil
}
