package main

import (
	"github.com/NubeDev/flow-framework/model"
)

//wizard make a network/dev/pnt
func (i *Instance) wizardSerial() (string, error) {
	var serial model.SerialConnection
	serial.SerialPort = "dev/ttyUSB0"
	serial.BaudRate = 9600

	var net model.Network
	net.Name = "modbus"
	net.TransportType = model.TransType.Serial
	net.PluginPath = "modbus"
	net.SerialConnection = &serial

	var dev model.Device
	dev.Name = "modbus"
	dev.AddressId = 1

	var pnt model.Point
	pnt.Name = "test"
	pnt.Description = "test"
	pnt.AddressId = 1
	pnt.ObjectType = "writeCoil"

	_, err = i.db.WizardNewNetDevPnt("modbus", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add modbus serial network wizard", err
	}

	return "pass: added network and points", err
}
