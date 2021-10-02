package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"time"
)

//pointUpdate update point present value
func (i *Instance) pointUpdate(uuid string, point *model.Point) (*model.Point, error) {
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	point.CommonFault.LastOk = time.Now().UTC()
	_, _ = i.db.UpdatePoint(uuid, point, true, true)
	if err != nil {
		log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
		return nil, err
	}
	return nil, nil
}

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
	dev.PollDelayPointsMS = 5000

	var pnt model.Point
	pnt.Name = "modbus"
	pnt.Description = "modbus"
	pnt.AddressId = utils.NewInt(1) //TODO check conversion
	pnt.ObjectType = model.ObjectTypes.WriteSingleFloat32

	_, err = i.db.WizardNewNetDevPnt("modbus", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add modbus TCP network wizard", err
	}

	return "pass: added network and points", err
}

//wizard make a network/dev/pnt
func (i *Instance) wizardSerial() (string, error) {
	var s model.SerialConnection
	s.SerialPort = "dev/ttyUSB0"
	s.BaudRate = 9600

	var net model.Network
	net.Name = "modbus"
	net.TransportType = model.TransType.Serial
	net.PluginPath = "modbus"
	net.SerialConnection = &s

	var dev model.Device
	dev.Name = "modbus"
	dev.AddressId = 1
	dev.ZeroMode = utils.NewTrue()
	dev.PollDelayPointsMS = 5000

	var pnt model.Point
	pnt.Name = "modbus"
	pnt.Description = "modbus"
	pnt.AddressId = utils.NewInt(1) //TODO check conversion
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
