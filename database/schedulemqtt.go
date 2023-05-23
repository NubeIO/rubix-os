package database

import (
	"fmt"
	"github.com/NubeIO/rubix-os/config"
	"github.com/NubeIO/rubix-os/services/localmqtt"
	"github.com/NubeIO/rubix-os/utils/boolean"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const fetchSchedulesTopic = "rubix/platform/schedules"

func (d *GormDatabase) PublishSchedulesListener() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		d.PublishSchedulesList(fmt.Sprintf("%s/publish", fetchSchedulesTopic))
	}
	topic := fetchSchedulesTopic
	mqttClient := localmqtt.GetLocalMqtt().Client
	if mqttClient != nil {
		err := mqttClient.Subscribe(topic, 1, callback)
		if err != nil {
			log.Errorf("localmqtt-broker subscribe: %s err: %s", topic, err.Error())
		} else {
			log.Infof("localmqtt-broker subscribe: %s", topic)
		}
	}
}

func (d *GormDatabase) PublishSchedulesList(topic string) {
	if boolean.IsFalse(config.Get().MQTT.PublishScheduleList) {
		return
	}
	data, err := d.GetSchedulesResult()
	if err != nil {
		log.Error("PublishSchedulesList error:", err)
		return
	}
	localmqtt.PublishSchedules(data, topic)
}
