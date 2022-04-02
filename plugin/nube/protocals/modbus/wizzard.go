package main

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

//wizard make a network/dev/pnt
func (inst *Instance) wizardTCP(body wizard) (string, error) {
	ip := "192.168.15.202"
	if body.IP != "" {
		ip = body.IP
	}
	p := 502
	if body.Port != 0 {
		p = body.Port
	}
	da := 1
	if body.DeviceAddr != 0 {
		da = int(body.BaudRate)
	}
	var net model.Network
	net.Name = "modbus"
	net.TransportType = model.TransType.IP
	net.PluginPath = "modbus"

	var dev model.Device
	dev.Name = "modbus"
	dev.CommonIP.Host = ip
	dev.CommonIP.Port = p
	dev.AddressId = da
	dev.ZeroMode = utils.NewTrue()
	dev.PollDelayPointsMS = 5000

	var pnt model.Point
	pnt.Name = "modbus"
	pnt.Description = "modbus"
	pnt.AddressID = utils.NewInt(1) //TODO check conversion
	pnt.ObjectType = string(model.ObjTypeWriteFloat32)

	_, err = inst.db.WizardNewNetworkDevicePoint("modbus", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add modbus TCP network wizard", err
	}

	return "pass: added network and points", err
}

//wizard make a network/dev/pnt
func (inst *Instance) wizardSerial(body wizard) (string, error) {

	sp := "/dev/ttyUSB0"
	if body.SerialPort != "" {
		sp = body.SerialPort
	}
	br := 9600
	if body.BaudRate != 0 {
		br = int(body.BaudRate)
	}

	var net model.Network
	net.Name = "modbus"
	net.TransportType = model.TransType.Serial
	net.PluginPath = "modbus"
	net.SerialPort = &sp
	net.SerialBaudRate = utils.NewUint(uint(br))
	net.SerialParity = utils.NewStr("none")
	net.SerialDataBits = utils.NewUint(1)
	net.SerialStopBits = utils.NewUint(1)

	da := 1
	if body.DeviceAddr != 0 {
		da = int(body.BaudRate)
	}

	var dev model.Device
	dev.Name = "modbus"
	dev.AddressId = da
	dev.ZeroMode = utils.NewTrue()
	dev.PollDelayPointsMS = 5000

	var pnt model.Point
	pnt.Name = "modbus"
	pnt.Description = "modbus"
	pnt.AddressID = utils.NewInt(1) //TODO check conversion
	pnt.ObjectType = string(model.ObjTypeWriteCoil)

	pntRet, err := inst.db.WizardNewNetworkDevicePoint("modbus", &net, &dev, &pnt)
	if err != nil {
		return "error: on flow-framework add modbus serial network wizard", err
	}

	log.Println(pntRet, err)
	return "pass: added network and points", err
}
