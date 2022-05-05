package main

import (
	"encoding/json"
	model "github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lwmodel"
	bm "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	elsysAPB = "ELSYS-ABP"
)

//mqttUpdate listen on mqtt and then update the point in flow-framework
func (inst *Instance) mqttUpdate(body mqtt.Message, devEUI, appID string) (*bm.Point, error) {
	//do an api call to chirpstack to get the device profile
	//decode the mqtt payload based of the device profile
	//if deviceProfileName

	payload := new(model.BasePayload)
	err := json.Unmarshal(body.Payload(), &payload)

	dev, err := inst.REST.GetDevice(payload.DevEUI)
	if err != nil {
		return nil, err
	}
	//check the payload for how to decode from
	if dev.Device.DeviceProfileID == elsysAPB {
		decoded := new(model.ElsysAPB)
		err = json.Unmarshal(body.Payload(), &decoded)
	}
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//addPoint from rest api
func (inst *Instance) addPoint(body *bm.Point) (*bm.Point, error) {
	return nil, nil
}

//pointPatch from rest
func (inst *Instance) pointPatch(body *bm.Point) (*bm.Point, error) {
	return nil, nil
}

//writePoint update point. Called via API call.
func (inst *Instance) writePoint(pntUUID string, body *bm.PointWriter) (point *bm.Point, err error) {
	//TODO: check for PointWriteByName calls that might not flow through the plugin.
	if body == nil {
		return
	}
	point, err = inst.db.WritePoint(pntUUID, body, true)
	if err != nil || point == nil {
		return nil, err
	}
	return point, nil
}

//delete point make sure
func (inst *Instance) deletePoint(body *bm.Point) (bool, error) {
	return true, nil
}

//DropDevices drop all devices
func (inst *Instance) DropDevices() (bool, error) {
	cli := inst.REST
	devices, err := cli.GetDevices()
	if err != nil {
		return false, err
	}
	for _, dev := range devices.Result {
		_, err := cli.DeleteDevice(dev.DevEUI)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

//delete point make sure
func (inst *Instance) serverDeletePoint(body *bm.Point) (bool, error) {
	return true, nil
}
