package main

import (
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/mqttclient"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
	plgrest "github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/restclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

//checkTypes make sure
func (i *Instance) checkTypes(body *model.Point) (*model.Point, error) {
	if body.ObjectType != model.ObjectTypeBACnet.AnalogValue || body.ObjectType != model.ObjectTypeBACnet.AnalogOutput {
		return nil, errors.New("data types supported are only AnalogValue or AnalogOutput")
	}
	return nil, nil
}

//bacnetUpdate listen on mqtt and then update the point in flow-framework
func (i *Instance) bacnetUpdate(body mqtt.Message) (*model.Point, error) {
	payload := new(pkgmodel.MqttPayload)
	err := json.Unmarshal(body.Payload(), &payload)
	t := mqttclient.TopicParts(body.Topic())
	top := t.Get(5)
	tt := top.(string)
	objType, addr := getPointAddr(tt)
	var point model.Point
	var pri model.Priority
	pri.P16 = payload.Value
	point.Priority = &pri
	pnt, _ := i.db.PointAndQuery(objType, addr) //TODO check if existing exists, as in the same addr
	if err != nil {
		log.Error("BACNET UPDATE POINT PointAndQuery")
		return nil, err
	}
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	point.CommonFault.LastOk = time.Now().UTC()
	_, _ = i.db.UpdatePoint(pnt.UUID, &point, false, true)
	if err != nil {
		log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
		return nil, err
	}
	return nil, nil
}

//addPoint from rest api
func (i *Instance) addPoint(body *model.Point) (*model.Point, error) {
	var point pkgmodel.BacnetPoint
	point.ObjectName = body.Name
	point.Enable = true
	point.Description = body.Description
	point.Address = body.AddressId
	point.ObjectType = body.ObjectType
	point.COV = body.COV
	point.EventState = "normal"
	point.Units = "noUnits"
	point.RelinquishDefault = body.Fallback
	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.AddPoint(point)
	//TODO check if existing exists, as in the same addr and also set the point in fault or out of fault
	if err != nil {
		return nil, err
	}
	return nil, nil

}

//pointPatch from rest
func (i *Instance) pointPatch(body *model.Point) (*model.Point, error) {
	var point pkgmodel.BacnetPoint
	point.Priority.P1 = body.Priority.P1
	point.Priority.P2 = body.Priority.P2
	point.Priority.P3 = body.Priority.P3
	point.Priority.P4 = body.Priority.P4
	point.Priority.P5 = body.Priority.P5
	point.Priority.P6 = body.Priority.P6
	point.Priority.P7 = body.Priority.P7
	point.Priority.P8 = body.Priority.P8
	point.Priority.P9 = body.Priority.P9
	point.Priority.P10 = body.Priority.P10
	point.Priority.P11 = body.Priority.P11
	point.Priority.P12 = body.Priority.P12
	point.Priority.P13 = body.Priority.P13
	point.Priority.P14 = body.Priority.P14
	point.Priority.P15 = body.Priority.P15
	point.Priority.P16 = body.Priority.P16
	point.ObjectName = body.Name
	addr := body.AddressId
	obj := body.ObjectType
	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.EditPoint(point, obj, addr)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

//delete point make sure
func (i *Instance) deletePoint(body *model.Point) (bool, error) {
	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.DeletePoint(body.ObjectType, body.AddressId)
	if err != nil {
		return false, err
	}
	return true, nil
}

//delete point make sure
func (i *Instance) bacnetServerDeletePoint(body *pkgmodel.BacnetPoint) (bool, error) {
	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.DeletePoint(body.ObjectType, body.Address)
	if err != nil {
		return false, err
	}
	return true, nil
}
