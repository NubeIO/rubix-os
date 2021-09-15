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
	if body.ObjectType != model.ObjectTypeBACnet.AnalogValue || body.ObjectType != model.ObjectTypeBACnet.AnalogOutput  {
		return nil, errors.New("data types supported are only AnalogValue or AnalogOutput")
	}
	return nil, nil
}



//checkTypes make sure
func (i *Instance) bacnetUpdate(body mqtt.Message) (*model.Point, error) {

	m := new(pkgmodel.MqttPayload)
	err := json.Unmarshal(body.Payload(), &m)
	if err != nil {

	}
	t := mqttclient.TopicParts(body.Topic())
	fmt.Println(t)
	top := t.Get(5)
	aaa := top.(string)
	objType, addr := getPointAddr(aaa)

	fmt.Println(objType, addr)
	var point pkgmodel.BacnetPoint
	_, _ = i.editPoint(point)
	if err != nil {
		return nil, err
	}
	return nil, nil
}


func (i *Instance) addPoint(body *model.Point) (*model.Point, error) {
	var point pkgmodel.BacnetPoint
	point.ObjectName = body.Name
	point.Address = body.AddressId
	point.ObjectType = body.ObjectType
	point.COV = body.COV
	point.EventState = "normal"
	point.Units = "noUnits"
	point.RelinquishDefault = 0.0
	fmt.Println(point.ObjectName, body.NextAvailableAddress, body.NextAvailableAddress)
	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.AddPoint(point)
	if err != nil {
		fmt.Println(err, 999999)
		return nil, err
	}
	fmt.Println(point.ObjectName, point.UseNextAvailableAddr)
	return nil, nil

}


func (i *Instance) editPoint(body pkgmodel.BacnetPoint) (*model.Point, error) {
	body.Priority.P1 = 1.0
	cli := plgrest.NewNoAuth(ip, port)
	_, err := cli.EditPoint(body)
	if err != nil {
		fmt.Println(err, 999999)
		return nil, err
	}

	return nil, nil

}