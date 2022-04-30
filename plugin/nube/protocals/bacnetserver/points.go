package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnetmodel"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/api/bacnetserver/v1/bsrest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/str"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

func setObjectType(objectType string) (objType string) {
	switch objectType {
	case string(model.ObjTypeAnalogInput), string(model.ObjTypeAnalogValue):
		objType = "analogValue"
	case string(model.ObjTypeAnalogOutput):
		objType = "analogOutput"
	case string(model.ObjTypeBinaryInput), string(model.ObjTypeBinaryValue):
		objType = "binaryValue"
	case string(model.ObjTypeBinaryOutput):
		objType = "binaryOutput"
	}
	return
}

//addPoint from rest api
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	bacnetPoint := &bsrest.BacnetPoint{}
	if body.Description == "" {
		bacnetPoint.Description = "na"
	}
	if nils.IntNilCheck(body.AddressID) || nils.IntIsNil(body.AddressID) == 0 {
		bacnetPoint.UseNextAvailableAddr = true
		//bacnetPoint.Address = nums.RandInt(1, 65000)
	} else {
		bacnetPoint.Address = nils.IntIsNil(body.AddressID)
	}
	if body.ObjectType == "" {
		bacnetPoint.ObjectType = "analog_value"
	} else {
		bacnetPoint.ObjectType = body.ObjectType
	}
	object := setObjectType(body.ObjectType)
	if object == "" {
		errMsg := fmt.Sprintf("bacnet-server: no object type passed in")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	bacnetPoint.ObjectType = object
	bacnetPoint.ObjectName = body.Name
	bacnetPoint.Enable = true

	bacnetPoint.COV = nils.Float64IsNil(body.COV) + 0.1
	bacnetPoint.EventState = "normal"
	bacnetPoint.Units = "noUnits"
	bacnetPoint.RelinquishDefault = nils.Float64IsNil(body.Fallback)
	log.Infoln("bacnet-server try and make a new point object type:", bacnetPoint.ObjectType, "object Address:", bacnetPoint.Address, "ObjectName:", bacnetPoint.ObjectName, "USE NEXT ADDRESS:", bacnetPoint.UseNextAvailableAddr)
	bacPoint, r := bacnetClient.AddPoint(bacnetPoint)
	err = r.GetError()
	if err != nil || bacPoint == nil {
		errMsg := fmt.Sprintf("bacnet-server.app.addPoint() add point was empty from bacnet-server")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	bacnetUUID := bacPoint.UUID
	body.AddressUUID = nils.NewString(bacnetUUID)
	body.AddressID = nils.NewInt(bacPoint.Address)
	point, err = inst.db.CreatePoint(body, true, false)
	if err != nil {
		if bacnetUUID == "" {
			errMsg := fmt.Sprintf("bacnet-server.app.addPoint() try and delete point over rest error: no point UUID was provided")
			log.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
		//if ail to add a new point in FF then delete it in the bacnet stack
		r, notFound, deleteOk := bacnetClient.DeletePoint(bacnetUUID)
		log.Infoln("bacnet-server.app.addPoint() delete the point if fail on add point in FF delete -> notFound:", notFound, "delete -> deleteOk:", deleteOk)
		if r.GetError() != nil {
			errMsg := fmt.Sprintf("bacnet-server: failed to add new point in bacnet stack, failed to remove the newly added point from bacnet-server-app error: %s", err.Error())
			log.Errorf(errMsg)
			return nil, errors.New(errMsg)
		} else {
			errMsg := fmt.Sprintf("bacnet-server: failed to add new point in flow-framwork, the point was deleted in bacnet-stack: %s", err.Error())
			log.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	return body, err
}

//updatePoint from rest
func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	bacnetPointUUID := nils.StringIsNil(body.AddressUUID)
	if bacnetPointUUID == "" || body == nil {
		log.Errorln("bacnet-server.app.updatePoint() body or address_uuid was empty")
		return nil, errors.New("no address_uuid")
	}
	bacnetPoint := &bsrest.BacnetPoint{}
	bacnetPoint.ObjectName = body.Name
	bacnetPoint.ObjectType = body.ObjectType
	if !utils.IntNilCheck(body.AddressID) {
		bacnetPoint.Address = utils.IntIsNil(body.AddressID)
	}
	log.Infoln("bacnet-server.app.updatePoint() try and update point over rest name:", body.Name, "bacnetPointUUID:", bacnetPointUUID)
	bacnetPoint, r := bacnetClient.UpdatePoint(bacnetPointUUID, bacnetPoint)
	err = errorMsg(r)
	if err != nil || bacnetPoint == nil {
		return nil, err
	}
	return body, err

}

//updatePointValue update point present value
func (inst *Instance) updatePointValue(body *model.Point) (*model.Point, error) {
	bacnetPointUUID := nils.StringIsNil(body.AddressUUID)
	if bacnetPointUUID == "" || body == nil {
		log.Errorln("bacnet-server.app.updatePointValue() body or address_uuid was empty")
		errMsg := fmt.Sprintf("bacnet-server: updatePointValue(): point is body is empty")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	bacnetPoint := &bsrest.BacnetPoint{}
	bacnetPoint.Priority = new(bsrest.Priority)
	if (*body.Priority).P16 != nil {
		(*bacnetPoint.Priority).P16 = (*body.Priority).P16
	}
	//if !utils.IntNilCheck(body.AddressID) {
	//	bacnetPoint.Address = utils.IntIsNil(body.AddressID)
	//}
	point, r := bacnetClient.UpdatePointValue(bacnetPointUUID, bacnetPoint)
	if r.GetError() != nil || point == nil {
		log.Errorln("bacnet-server: updatePointValue() body back from rest was nil or err:", r.GetError())
	} else {
		log.Println("bacnet-server: updatePointValue() point-name:", point.ObjectName)
	}
	return nil, nil
}

//deletePoint point make sure
func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	if nils.StringNilCheck(body.AddressUUID) { //still delete it anyway
		log.Errorln("bacnet-server.app.deletePoint() no address_uuid provided")
		ok, err = inst.db.DeletePoint(body.UUID)
		if err != nil {
			return false, err
		}
	} else {
		ok, err = inst.db.DeletePoint(body.UUID) //delete and try and delete on bacnet-server
		if err != nil {
			return false, err
		}
		r, notFound, deletedOk := bacnetClient.DeletePoint(nils.StringIsNil(body.AddressUUID))
		log.Infoln("bacnet-server.app.deletePoint() statusCode:", r.GetStatusCode(), "notFound", notFound, "deletedOk", deletedOk)
	}
	return

}

func (inst *Instance) dropPoints() (err error) {
	r := bacnetClient.DropPoints()
	err = r.GetError()
	return
}

//bacnetUpdate listen on mqtt and then update the point in flow-framework
func (inst *Instance) bacnetUpdate(body mqtt.Message) (*model.Point, error) {
	payload := new(bacnetmodel.MqttPayload)
	err := json.Unmarshal(body.Payload(), &payload)
	t, _ := mqttclient.TopicParts(body.Topic())
	const objectType = 10
	top := t.Get(objectType)
	tt := top.(string)
	objType, addr := getPointAddr(tt)
	var point model.Point
	var pri model.Priority
	pri.P16 = payload.Value
	point.Priority = &pri
	addressID := addr

	object := str.NewString(objType).ToSnakeCase()
	object = str.NewString(object).LcFirstLetter()
	pnt, err := inst.db.GetOnePointByArgs(api.Args{ObjectType: nils.NewString(object), AddressID: nils.NewString(addressID)})
	if err != nil {
		log.Error("bacnet-server: GetOnePointByArgs() issue on message from mqtt update point", err)
		return nil, err
	}
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	point.CommonFault.LastOk = time.Now().UTC()
	if pnt == nil {
		log.Error("bacnet-server: issue on message from mqtt update point")
		return nil, err
	}
	_, err = inst.db.UpdatePointValue(pnt.UUID, &point, true)
	if err != nil {
		log.Error("bacnet-server: UpdatePointValue() issue on message from mqtt update point")
		return nil, err
	}
	return nil, nil
}
