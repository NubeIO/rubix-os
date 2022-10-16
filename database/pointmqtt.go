package database

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/services/mqtt"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) PublishPointsList() {
	networks, err := d.GetNetworks(api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		log.Error("PublishPointsList error:", err)
		return
	}
	mqtt.PublishPointsList(networks)
}

func (d *GormDatabase) RePublishPointsCov() {
	networks, err := d.GetNetworks(api.Args{WithDevices: true, WithPoints: true, WithPriority: true})
	if err != nil {
		log.Error("RePublishPointsCov error:", err)
		return
	}
	for _, network := range networks {
		for _, device := range network.Devices {
			for _, point := range device.Points {
				priority := point.Priority.GetHighestPriorityValue()
				mqtt.PublishPointCov(network, device, point, priority)
			}
		}
	}
}

func (d *GormDatabase) PublishPointCov(uuid string) error {
	point, err := d.GetPoint(uuid, api.Args{WithPriority: true})
	if err != nil {
		return err
	}
	priority := point.Priority.GetHighestPriorityValue()
	device, err := d.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil {
		return err
	}
	network, err := d.GetNetwork(device.NetworkUUID, api.Args{})
	if err != nil {
		return err
	}
	mqtt.PublishPointCov(network, device, point, priority)
	return nil
}
