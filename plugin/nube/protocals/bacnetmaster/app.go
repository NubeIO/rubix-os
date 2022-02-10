package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nrest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nums"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/system/networking"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) addNetwork(body *model.Network) (*Network, error) {
	newNetwork := new(Network)
	if body.NetworkInterface != "" {
		_net, _ := networking.GetInterfaceByName(body.NetworkInterface)
		if _net == nil {
			log.Error("bacnet-master-plugin: ERROR failed to find a valid network interface")
			return nil, errors.New("failed to find a valid network interface")
		}
		newNetwork.NetworkIp = _net.IP
		newNetwork.NetworkMask = _net.NetMaskLength
	}
	newNetwork.NetworkName = body.Name
	rt.Method = nrest.PUT
	rt.Path = networkBacnet
	res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: newNetwork})
	if err != nil {
		log.Error("bacnet-master-plugin: ERROR added network over rest response-code:", code)
		return nil, nil
	}

	res.ToInterfaceNoErr(newNetwork)
	body.AddressUUID = newNetwork.NetworkUUID
	log.Println("bacnet-master-plugin: added network over rest response-code:", code, "bacnet-master network uuid", newNetwork.NetworkUUID)
	updateNetwork, err := i.db.UpdateNetwork(body.UUID, body)
	if err != nil {
		log.Error("bacnet-master-plugin: ERROR update network addressUUID", updateNetwork.UUID)
		return nil, err
	}
	log.Println("bacnet-master-plugin: update network addressUUID", updateNetwork.UUID)

	return newNetwork, nil
}

//addDevice from rest api
func (i *Instance) addDevice(body *model.Device) (*model.Device, error) {
	newDevice := new(Device)
	newDevice.DeviceName = body.Name
	newDevice.DeviceIp = body.CommonIP.Host
	newDevice.DevicePort = body.CommonIP.Port
	newDevice.DeviceMask = nums.IntIsNil(body.DeviceMask)
	newDevice.DeviceObjectId = 1

	getNetwork, err := i.db.GetNetwork(body.NetworkUUID, api.Args{})
	if err != nil {
		log.Error("bacnet-master-plugin: ERROR get network addressUUID", body.NetworkUUID)
		return nil, err
	}
	if getNetwork == nil {
		log.Error("bacnet-master-plugin: ERROR failed to find a network", body.NetworkUUID)
		return nil, errors.New("bacnet-master-plugin: ERROR failed to find a network")
	}
	fmt.Println(getNetwork)
	fmt.Println(12222222, getNetwork.UUID, getNetwork.AddressUUID)
	newDevice.NetworkUuid = getNetwork.AddressUUID
	rt.Method = nrest.PUT
	rt.Path = deviceBacnet
	res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: newDevice})
	if err != nil {
		log.Error("bacnet-master-plugin: ERROR added device over rest response-code:", code, res.AsString())
		return nil, nil
	}

	res.ToInterfaceNoErr(newDevice)
	body.AddressUUID = newDevice.NetworkUuid
	log.Println("bacnet-master-plugin: added device over rest response-code:", code, "bacnet-master device uuid", newDevice.DeviceUUID)
	updateDev, err := i.db.UpdateDevice(body.UUID, body)
	if err != nil {
		log.Error("bacnet-master-plugin: ERROR update device addressUUID", updateDev.UUID)
		return nil, err
	}
	log.Println("bacnet-master-plugin: update device addressUUID", updateDev.UUID)
	return nil, nil

}

//addPoint from rest api
func (i *Instance) addPoint(body *model.Point) (*model.Point, error) {
	newPoint := new(Point)
	newPoint.PointName = body.Name
	newPoint.PointObjectId = nums.IntIsNil(body.AddressID)
	newPoint.PointObjectType = body.ObjectType

	getDevice, err := i.db.GetDevice(body.DeviceUUID, api.Args{})
	if err != nil {
		log.Error("bacnet-master-plugin: ERROR get device addressUUID", body.DeviceUUID)
		return nil, err
	}
	if getDevice == nil {
		log.Error("bacnet-master-plugin: ERROR failed to find a device", body.DeviceUUID)
		return nil, errors.New("bacnet-master-plugin: ERROR failed to find a device")
	}

	newPoint.DeviceUuid = getDevice.AddressUUID

	rt.Method = nrest.PUT
	rt.Path = pointBacnet
	res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{Json: newPoint})
	if err != nil {
		log.Error("bacnet-master-plugin: ERROR added point over rest response-code:", code)
		return nil, nil
	}

	res.ToInterfaceNoErr(newPoint)
	body.AddressUUID = newPoint.DeviceUuid
	log.Println("bacnet-master-plugin: added device over rest response-code:", code, "bacnet-master point uuid", newPoint.DeviceUuid)
	updateDev, err := i.db.UpdatePoint(body.UUID, body, true)
	if err != nil {
		log.Error("bacnet-master-plugin: ERROR update point addressUUID", updateDev.UUID)
		return nil, err
	}
	log.Println("bacnet-master-plugin: update point addressUUID", updateDev.UUID)
	return nil, nil

}

//bacnetUpdate listen on mqtt and then update the point in flow-framework
func (i *Instance) bacnetUpdate(body mqtt.Message) (*model.Point, error) {
	payload := new(bacnet_model.MqttPayload)
	err := json.Unmarshal(body.Payload(), &payload)
	t, _ := mqttclient.TopicParts(body.Topic())
	const objectType = 10
	top := t.Get(objectType)
	tt := top.(string)
	objType, addr := getPointAddr(tt)

	fmt.Println(err, objType, addr)
	return nil, nil

}

//pointPatch from rest
func (i *Instance) pointPatch(body *model.Point) (*model.Point, error) {

	return nil, nil
}

//deletePoint point make sure
func (i *Instance) deletePoint(body *model.Point) (bool, error) {

	return true, nil
}

//bacnetServerDeletePoint point make sure
func (i *Instance) bacnetServerDeletePoint(body *bacnet_model.BacnetPoint) (bool, error) {
	//cli := plgrest.NewNoAuth(ip, port)
	//_, err := cli.DeletePoint(body.ObjectType, body.Address)
	//if err != nil {
	//	return false, err
	//}
	return true, nil
}
