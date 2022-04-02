package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/api/bacnetserver/v1/bsrest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/api/common/v1/iorest"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

//addNetwork add network
func (inst *Instance) addNetwork(body *model.Network) (network *model.Network, err error) {
	nets, err := inst.db.GetNetworksByPluginName(body.PluginPath, api.Args{})
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
	network, err = inst.db.CreateNetwork(body, true)
	if err != nil {
		return nil, err
	}
	return body, nil
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

//addPoint from rest api
func (inst *Instance) addPoint(body *model.Point) (point *model.Point, err error) {
	bacnetPoint := bsrest.BacnetPoint{}
	if body.Description == "" {
		bacnetPoint.Description = "na"
	}
	bacnetPoint.ObjectName = body.Name
	bacnetPoint.Enable = true
	bacnetPoint.Address = utils.IntIsNil(body.AddressID)
	bacnetPoint.ObjectType = body.ObjectType
	bacnetPoint.COV = utils.Float64IsNil(body.COV)
	bacnetPoint.EventState = "normal"
	bacnetPoint.Units = "noUnits"
	bacnetPoint.RelinquishDefault = utils.Float64IsNil(body.Fallback)

	bacPoint, r := bacnetClient.AddPoint(bacnetPoint)
	err = r.GetError()
	if err != nil {
		return nil, err
	}
	bacnetUUID := bacPoint.UUID
	body.AddressUUID = nils.NewString(bacnetUUID)
	point, err = inst.db.CreatePoint(body, true, false)
	if err != nil {
		//if ail to add a new point in FF then delete it in the bacnet stack
		r, notFound, deleteOk := bacnetClient.DeletePoint(bacnetUUID)
		fmt.Println(notFound, deleteOk)
		if r.GetError() != nil {
			errMsg := fmt.Sprintf("bacnet-server: failed to add new point in bacnet stack, failed to remove the newly added point from bacnet-server-app error: %s", err.Error())
			log.Errorf(errMsg)
			return nil, err
		} else {
			errMsg := fmt.Sprintf("bacnet-server: failed to add new point in flow-framwork, the point was deleted in bacnet-stack: %s", err.Error())
			log.Errorf(errMsg)
			return nil, err
		}
	}
	return body, err
}

//updateNetwork update network
func (inst *Instance) updateNetwork(body *model.Network) (network *model.Network, err error) {
	network, err = inst.db.UpdateNetwork(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return network, nil
}

//updateDevice update device
func (inst *Instance) updateDevice(body *model.Device) (device *model.Device, err error) {
	device, err = inst.db.UpdateDevice(body.UUID, body, true)
	if err != nil {
		return nil, err
	}
	return device, nil
}

//updatePoint from rest
func (inst *Instance) updatePoint(body *model.Point) (point *model.Point, err error) {
	bacnetPointUUID := nils.StringIsNil(body.AddressUUID)
	if bacnetPointUUID == "" {
		return nil, errors.New("no address_uuid")
	}
	bacnetPoint := bsrest.BacnetPoint{}
	bacnetPoint.ObjectName = body.Name
	bacnetPoint.ObjectType = body.ObjectType
	if !utils.IntNilCheck(body.AddressID) {
		bacnetPoint.Address = utils.IntIsNil(body.AddressID)
	}
	bacnetPoint, r := bacnetClient.UpdatePoint(bacnetPointUUID, bacnetPoint)
	err = errorMsg(r)
	if err != nil {
		return nil, err
	}
	return body, err

}

//updatePointValue update point present value
func (inst *Instance) updatePointValue(body *model.Point) (*model.Point, error) {

	bacnetPoint := bsrest.BacnetPoint{}
	bacnetPoint.Priority = new(bsrest.Priority)
	if (*body.Priority).P16 != nil {
		(*bacnetPoint.Priority).P16 = (*body.Priority).P16
	}
	bacnetPointUUID := nils.StringIsNil(body.AddressUUID)
	if !utils.IntNilCheck(body.AddressID) {
		bacnetPoint.Address = utils.IntIsNil(body.AddressID)
	}

	point, r := bacnetClient.UpdatePointValue(bacnetPointUUID, bacnetPoint)

	if r.GetError() != nil {
		log.Errorln("bacnet-server: updatePointValue() point-name:", point.ObjectName)
		log.Errorln("bacnet-server: updatePointValue() err:", r.GetError())
	} else {
		log.Println("bacnet-server: updatePointValue() point-name:", point.ObjectName)
	}

	return nil, nil
}

//deleteNetwork network
func (inst *Instance) deleteNetwork(body *model.Network) (ok bool, err error) {
	err = inst.dropPoints()
	if err != nil {
		return ok, err
	}
	ok, err = inst.db.DeleteNetwork(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//deleteNetwork device
func (inst *Instance) deleteDevice(body *model.Device) (ok bool, err error) {
	err = inst.dropPoints()
	if err != nil {
		return ok, err
	}
	ok, err = inst.db.DeleteDevice(body.UUID)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//deletePoint point make sure
func (inst *Instance) deletePoint(body *model.Point) (ok bool, err error) {
	if body.AddressUUID == nil {
		return false, errors.New("no provided address_uuid")
	}
	r, notFound, deletedOk := bacnetClient.DeletePoint(nils.StringIsNil(body.AddressUUID))
	//if point not found lets still delete from FF
	if notFound || deletedOk {
		ok, err = inst.db.DeletePoint(body.UUID)
		if err != nil {
			return false, err
		}
		ok = true
		return ok, nil
	} else {
		err = r.GetError()
		return ok, err
	}
}

func (inst *Instance) dropPoints() (err error) {
	r := bacnetClient.DropPoints()
	err = r.GetError()
	return
}

//

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
	pnt, err := inst.db.GetOnePointByArgs(api.Args{ObjectType: &objType, AddressID: &addr})
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

func errorMsg(response *iorest.RestResponse) (err error) {
	msg := response.Message
	if response.BadRequest {
		err = fmt.Errorf("%s:  msg:%s", "bacnet-server", msg)
	}
	return

}
