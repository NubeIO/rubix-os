package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nrest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube_api"
	nube_api_bacnetserver "github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube_api/bacnetserver"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

//addNetwork add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	nets, err := inst.db.GetNetworksByName(body.PluginPath, api.Args{})
	if err != nil {
		return nil, err
	}
	for _, net := range nets {
		if net != nil {
			errMsg := fmt.Sprintf("bacnet-server: only max one network is allowed with bacnet-server")
			log.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	network, err = inst.db.CreateNetwork(body, false)
	if err != nil {
		return nil, err
	}
	return network, nil
}

//deleteNetwork network
func (inst *Instance) deleteNetwork() (ok bool, err error) {
	err = inst.dropPoints()
	if err != nil {
		return ok, err
	}
	ok = true
	return ok, nil
}

//deleteNetwork device
func (inst *Instance) deleteDevice() (ok bool, err error) {
	err = inst.dropPoints()
	if err != nil {
		return ok, err
	}
	ok = true
	return ok, nil
}

func (inst *Instance) dropPoints() (err error) {
	r := bacnetClient.DropPoints()
	err = errorMsg(r.Response)
	return
}

//addDevice add device
func (inst *Instance) addDevice(body *model.Device) (device *model.Device, err error) {
	network, err := inst.db.GetNetwork(body.NetworkUUID, api.Args{WithDevices: true})
	if err != nil {
		return nil, err
	}
	if len(network.Devices) >= 1 {
		errMsg := fmt.Sprintf("bacnet-server: only max one device is allowed with bacnet-server")
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	device, err = inst.db.CreateDevice(body)
	if err != nil {
		return nil, err
	}
	return device, nil
}

//bacnetUpdate listen on mqtt and then update the point in flow-framework
func (inst *Instance) bacnetUpdate(body mqtt.Message) (*model.Point, error) {
	payload := new(bacnet_model.MqttPayload)
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
	pnt, _ := inst.db.GetOnePointByArgs(api.Args{ObjectType: &objType, AddressID: &addr}) //TODO check conversion if existing exists, as in the same addr
	if err != nil {
		log.Error("BACNET UPDATE POINT PointAndQuery")
		return nil, err
	}
	point.CommonFault.InFault = false
	point.CommonFault.MessageLevel = model.MessageLevel.Info
	point.CommonFault.MessageCode = model.CommonFaultCode.Ok
	point.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	point.CommonFault.LastOk = time.Now().UTC()
	if pnt == nil {
		log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
		return nil, err
	}
	_, err = inst.db.UpdatePointValue(pnt.UUID, &point, true)
	if err != nil {
		log.Error("BACNET UPDATE POINT issue on message from mqtt update point")
		return nil, err
	}
	return nil, nil
}

func errorMsg(response nube_api.Response) (err error) {
	appName := reqType.Path
	msg := response.Message
	if response.BadRequest {
		err = fmt.Errorf("%s:  msg:%s", appName, msg)
	}
	return

}

//addPoint from rest api
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	bacPoint := nube_api_bacnetserver.BacnetPoint{}
	if body.Description == "" {
		bacPoint.Description = "na"
	}
	bacPoint.ObjectName = body.Name
	bacPoint.Enable = true
	bacPoint.Address = utils.IntIsNil(body.AddressID)
	bacPoint.ObjectType = body.ObjectType
	bacPoint.COV = utils.Float64IsNil(body.COV)
	bacPoint.EventState = "normal"
	bacPoint.Units = "noUnits"
	bacPoint.RelinquishDefault = utils.Float64IsNil(body.Fallback)

	bacPoint, r := bacnetClient.AddPoint(bacPoint)
	err = errorMsg(r.Response)
	if err != nil {
		return nil, err
	}
	bacnetUUID := bacPoint.UUID
	body.AddressUUID = nils.NewString(bacnetUUID)
	point, err = inst.db.CreatePoint(body, true, false)
	if err != nil {
		//if ail to add a new point in FF then delete it in the bacnet stack
		bacnetClient.DeletePoint(bacnetUUID)
		err = errorMsg(r.Response)
		if err != nil {
			errMsg := fmt.Sprintf("bacnet-server: failed to add new point in bacnet stack, failed to remove the newly added point from bacnet-server-app error: %s", err.Error())
			log.Errorf(errMsg)
			return nil, err
		}
	}
	return body, err
}

//pointPatch from rest
func (inst *Instance) pointPatch2(body *model.Point) (bacnetPoint nube_api_bacnetserver.BacnetPoint, err error) {
	//if body.AddressUUID == nil {
	//	return nil, errors.New("no address_uuid")
	//}

	//bacnetPoint.Priority = new(nube_api_bacnetserver.Priority)
	//if (*body.Priority).P16 != nil {
	//	(*bacnetPoint.Priority).P16 = (*body.Priority).P16
	//}
	//bacnetPoint.ObjectName = body.Name
	//bacnetPoint.Address = utils.IntIsNil(body.AddressID)
	//bacnetPoint.ObjectType = body.ObjectType
	//point.Units = body.Unit

	if !utils.IntNilCheck(body.AddressID) {
		//bacnetPoint.Address = utils.IntIsNil(body.AddressID)
		bacnetPoint.UseNextAvailableAddr = false
	}

	bacnetPoint.Description = body.Description

	return bacnetPoint, nil

}

//pointPatch from rest
func (inst *Instance) pointPatch(body *model.Point) (*model.Point, error) {
	if body.AddressUUID == nil {
		return nil, errors.New("no address_uuid")
	}
	point := new(bacnet_model.BacnetPoint)
	point.Priority = new(model.Priority)
	if (*body.Priority).P16 != nil {
		(*point.Priority).P16 = (*body.Priority).P16
	}
	point.ObjectName = body.Name
	point.Address = utils.IntIsNil(body.AddressID)
	point.ObjectType = body.ObjectType
	//point.Units = body.Unit
	point.Description = body.Description

	rt.Method = nrest.PATCH
	rt.Path = fmt.Sprintf("/api/bacnet/points/uuid/%s", *body.AddressUUID)
	_, _, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: point})

	if err != nil {
		log.Errorf("BACNET: EDIT POINT issue on add rest: %v\n", err)
		return nil, err
	}
	return nil, nil

}

//deletePoint point make sure
func (inst *Instance) deletePoint(body *model.Point) (bool, error) {
	if body.AddressUUID == nil {
		return false, errors.New("no address_uuid")
	}
	rt.Method = nrest.DELETE
	rt.Path = fmt.Sprintf("/api/bacnet/points/uuid/%s", *body.AddressUUID)
	_, _, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
	if err != nil {
		log.Errorf("BACNET: DELETE POINT issue on add rest: %v\n", err)
		return false, err
	}
	return true, nil
}

//bacnetServerDeletePoint point make sure
func (inst *Instance) bacnetServerDeletePoint(body *bacnet_model.BacnetPoint) (bool, error) {
	//cli := plgrest.NewNoAuth(ip, port)
	//_, err := cli.DeletePoint(body.ObjectType, body.Address)
	//if err != nil {
	//	return false, err
	//}
	return true, nil
}
