package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/mqttclient"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
	plgrest "github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/restclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	if err != nil {

	}
	t := mqttclient.TopicParts(body.Topic())
	top := t.Get(5)
	aaa := top.(string)
	objType, addr := getPointAddr(aaa)
	var point model.Point
	var pri model.Priority
	pri.P1 = payload.Value
	point.Priority = &pri
	pnt, _ := i.db.PointAndQuery(objType, addr)

	//TODO check if existing exists, as in the same addr

	if err != nil {
		return nil, err
	}
	_, _ = i.db.UpdatePoint(pnt.UUID, &point, false, true)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//addPoint from rest api
func (i *Instance) addPoint(body *model.Point) (*model.Point, error) {
	var point pkgmodel.BacnetPoint
	point.ObjectName = body.Name
	point.Address = body.AddressId
	point.ObjectType = body.ObjectType
	point.COV = body.COV
	point.EventState = "normal"
	point.Units = "noUnits"
	point.RelinquishDefault = 0.0
	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.AddPoint(point)
	//TODO check if existing exists, as in the same addr and also set the point in fault or out of fault
	if err != nil {
		return nil, err
	}
	fmt.Println(point.ObjectName, point.UseNextAvailableAddr)
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
