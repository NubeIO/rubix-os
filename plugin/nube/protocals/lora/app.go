package main

import (
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/tty"
	log "github.com/sirupsen/logrus"
)

func SerialOpenAndRead() error {
	s := new(tty.SerialSetting)
	s.BaudRate = 38400
	sc := tty.New(s)
	sc.NewSerialConnection()
	sc.Loop()
	return nil
}

// SerialOpen open serial port
func (c *Instance) SerialOpen() error {
	go func() error {
		err := SerialOpenAndRead();if err != nil {
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




