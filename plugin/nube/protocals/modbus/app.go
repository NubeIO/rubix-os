package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"go.bug.st/serial"
)

//wizard make a network/dev/pnt
func (i *Instance) wizardTCP() (string, error) {

	var net model.Network
	net.Name = "modbus"
	net.TransportType = model.TransType.IP
	net.PluginPath = "modbus"

	var dev model.Device
	dev.Name = "modbus"
	dev.CommonIP.Host = "192.168.15.202"
	dev.CommonIP.Port = 502
	dev.AddressId = 1
	dev.ZeroMode = utils.NewTrue()

	var pnt model.Point
	pnt.Name = "modbus"
	pnt.Description = "modbus"
	pnt.AddressId = 1
	pnt.ObjectType = model.ObjectTypes.WriteSingleFloat32

	_, err = i.db.WizardNewNetDevPnt("modbus", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add modbus TCP network wizard", err
	}

	return "pass: added network and points", err
}

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
	dev.ZeroMode = utils.NewTrue()

	var pnt model.Point
	pnt.Name = "modbus"
	pnt.Description = "modbus"
	pnt.AddressId = 1
	pnt.ObjectType = model.ObjectTypes.WriteCoil

	_, err = i.db.WizardNewNetDevPnt("modbus", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add modbus serial network wizard", err
	}

	return "pass: added network and points", err
}

//listSerialPorts list all serial ports on host
func (i *Instance) listSerialPorts() (*utils.Array, error) {
	ports, err := serial.GetPortsList()
	p := utils.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}
