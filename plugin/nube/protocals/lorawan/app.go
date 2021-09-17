package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//mqttUpdate listen on mqtt and then update the point in flow-framework
func (i *Instance) mqttUpdate(body mqtt.Message) (*model.Point, error) {
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

//delete point make sure
func (i *Instance) deletePoint(body *model.Point) (bool, error) {
	return true, nil
}

//DropDevices drop all devices
func (i *Instance) DropDevices() (bool, error) {
	cli := i.CLI
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
func (i *Instance) serverDeletePoint(body *pkgmodel.BacnetPoint) (bool, error) {
	return true, nil
}
