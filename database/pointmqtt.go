package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/services/localmqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const fetchDeviceInfo = "rubix/platform/info"
const fetchPointsTopic = "rubix/platform/points"

func (d *GormDatabase) PublishPointsListListener() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		d.PublishPointsList(fmt.Sprintf("%s/publish", fetchPointsTopic))
	}
	topic := fetchPointsTopic
	mqttClient := localmqtt.GetPointMqtt().Client
	if mqttClient != nil {
		err := mqttClient.Subscribe(topic, 1, callback)
		if err != nil {
			log.Errorf("localmqtt-broker subscribe:%s err:%s", topic, err.Error())
		} else {
			log.Infof("localmqtt-broker subscribe:%s", topic)
		}
	}
}

func (d *GormDatabase) PublishDeviceInfo() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		localmqtt.PublishInfo()
	}
	topic := fetchDeviceInfo
	mqttClient := localmqtt.GetPointMqtt().Client
	if mqttClient != nil {
		err := mqttClient.Subscribe(fetchDeviceInfo, 1, callback)
		if err != nil {
			log.Errorf("localmqtt-broker subscribe:%s err:%s", topic, err.Error())
		} else {
			log.Infof("localmqtt-broker subscribe:%s", topic)
		}
	}
}

func (d *GormDatabase) PublishPointsList(topic string) {
	networks, err := d.GetNetworks(api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		log.Error("PublishPointsList error:", err)
		return
	}
	localmqtt.PublishPointsList(networks, topic)
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
				localmqtt.PublishPointCov(network, device, point, priority)
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
	localmqtt.PublishPointCov(network, device, point, priority)
	return nil
}
