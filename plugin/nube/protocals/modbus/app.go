package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/utilstime"
	"time"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

//addDevice add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = inst.db.CreateNetwork(body, false)
	if err != nil {
		return nil, err
	}
	return network, nil
}

//addDevice add device
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	return device, nil
}

//addPoint add point
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	point, err = inst.db.CreatePoint(body, false, false)
	if err != nil {
		return nil, err
	}
	return point, nil
}

//pointUpdate update point present value
func (inst *Instance) pointUpdate(uuid string, value float64) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = fmt.Sprintf("last-update: %s", utilstime.TimeStamp())
	point.CommonFault.LastOk = time.Now().UTC()
	var pri model.Priority
	pri.P16 = &value
	point.Priority = &pri
	point.InSync = utils.NewTrue()
	_, err = inst.db.UpdatePointValue(uuid, &point, true)
	if err != nil {
		log.Error("MODBUS UPDATE POINT UpdatePointValue()", err)
		return nil, err
	}
	return nil, nil
}

//pointUpdate update point present value
func (inst *Instance) pointUpdateErr(uuid string, err error) (*model.Point, error) {
	var point model.Point
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	_, err = inst.db.UpdatePoint(uuid, &point, true)
	if err != nil {
		log.Error("MODBUS UPDATE POINT pointUpdateErr()", err)
		return nil, err
	}
	return nil, nil
}

//listSerialPorts list all serial ports on host
func (inst *Instance) listSerialPorts() (*utils.Array, error) {
	ports, err := serial.GetPortsList()
	p := utils.NewArray()
	for _, port := range ports {
		p.Add(port)
	}
	return p, err
}
