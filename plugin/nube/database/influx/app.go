package main

import (
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/bacnetserver/model"
)

const (
	elsysAPB = "ELSYS-ABP"
)

//DropDevices drop all devices
func (i *Instance) DropDevices() (bool, error) {
	cli := i.REST
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
