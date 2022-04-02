package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (d *GormDatabase) deviceNameExists(dev *model.Device, body *model.Device) bool {
	var arg api.Args
	arg.WithDevices = true
	device, err := d.GetNetwork(dev.NetworkUUID, arg)
	if err != nil {
		return false
	}
	for _, p := range device.Devices {
		if p.Name == body.Name {
			if p.UUID == dev.UUID {
				return false
			} else {
				return true
			}
		}
	}
	return false
}
