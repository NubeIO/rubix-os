package database

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/services/localmqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const fetchDeviceInfo = "rubix/platform/info"
const fetchPointsTopic = "rubix/platform/points"
const fetchPointTopicWrite = "rubix/platform/point/write"
const fetchPointTopic = "rubix/platform/point"

func (d *GormDatabase) PublishPointWriteListener() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		body := &interfaces.MqttPoint{}
		err := json.Unmarshal(message.Payload(), &body)
		if err == nil {
			if body != nil {
				d.publishPointWrite(body)
			}
		}
	}
	topic := fetchPointTopicWrite
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

func (d *GormDatabase) publishPointWrite(details *interfaces.MqttPoint) {
	if details == nil {
		return
	}
	if details.PointUUID != "" {
		_, _, _, _, err := d.PointWrite(details.PointUUID, details.Priority, false, false, nil)
		if err != nil {
			log.Error("mqtt write point by uuid: error:", err)
			return
		}
	} else {
		networkName := details.NetworkName
		deviceName := details.DeviceName
		pointName := details.PointName
		_, err := d.PointWriteByName(networkName, deviceName, pointName, details.Priority, false)
		if err != nil {
			log.Error("mqtt write point by name: error:", err)
			return
		}
	}
}

func (d *GormDatabase) PublishFetchPointListener() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		body := &interfaces.MqttPoint{}
		err := json.Unmarshal(message.Payload(), &body)
		if err == nil {
			if body != nil {
				d.PublishPoint(body)
			}
		}
	}
	topic := fetchPointTopic
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

func (d *GormDatabase) PublishPoint(details *interfaces.MqttPoint) {
	if details.PointUUID != "" {
		point, err := d.GetPoint(details.PointUUID, api.Args{WithPriority: true})
		if err != nil {
			log.Error("PublishPointsList error:", err)
			return
		}
		localmqtt.PublishPoint(point)
	} else {
		networks, err := d.GetNetworks(api.Args{WithDevices: true, WithPoints: true})
		if err != nil {
			log.Error("PublishPointsList error:", err)
			return
		}
		localmqtt.PublishPointByName(networks, details)
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
				localmqtt.PublishPointCov(network, device, point)
			}
		}
	}
}

func (d *GormDatabase) PublishPointCov(uuid string) error {
	point, err := d.GetPoint(uuid, api.Args{WithPriority: true})
	if err != nil {
		return err
	}
	device, err := d.GetDevice(point.DeviceUUID, api.Args{})
	if err != nil {
		return err
	}
	network, err := d.GetNetwork(device.NetworkUUID, api.Args{})
	if err != nil {
		return err
	}
	localmqtt.PublishPointCov(network, device, point)
	return nil
}
