package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/bacnetserver/bacnet_model"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nrest"
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
	body.NetworkAddressUUID = newNetwork.NetworkUUID
	log.Println("bacnet-master-plugin: added network over rest response-code:", code, "bacnet-master network uuid", newNetwork.NetworkUUID)
	updateNetwork, err := i.db.UpdateNetwork(body.UUID, body)
	if err != nil {
		log.Error("bacnet-master-plugin: ERROR update network addressUUID", updateNetwork.UUID)
		return nil, err
	}
	log.Println("bacnet-master-plugin: update network addressUUID", updateNetwork.UUID)

	return newNetwork, nil
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

//addPoint from rest api
func (i *Instance) addDevice(body *model.Device) (*model.Point, error) {

	return nil, nil

}

//addPoint from rest api
func (i *Instance) addPoint(body *model.Point) (*model.Point, error) {

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
