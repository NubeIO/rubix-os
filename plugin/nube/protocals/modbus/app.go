package main

import (
	"fmt"
	"time"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

//pointUpdate update point present value
func (i *Instance) pointUpdate(point *model.Point, value float64) (*model.Point, error) {
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	point.CommonFault.LastOk = time.Now().UTC()
	var pri model.Priority
	pri.P16 = &value
	point.Priority = &pri
	fmt.Println()
	_, _ = i.db.UpdatePointValue(point.UUID, point, true)
	if err != nil {
		log.Error("MODBUS UPDATE POINT issue on message from mqtt update point", err)
		return nil, err
	}
	return nil, nil
}

//pointUpdate update point present value
func (i *Instance) pointUpdateErr(uuid string, point *model.Point, err error) (*model.Point, error) {
	point.CommonFault.InFault = true
	point.CommonFault.MessageLevel = model.MessageLevel.Fail
	point.CommonFault.MessageCode = model.CommonFaultCode.PointError
	point.CommonFault.Message = err.Error()
	point.CommonFault.LastFail = time.Now().UTC()
	_, _ = i.db.UpdatePoint(uuid, point, true)
	if err != nil {
		log.Error("MODBUS UPDATE POINT issue on message from mqtt update point", err)
		return nil, err
	}
	return nil, nil
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
