package database

import (
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/services/localmqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (d *GormDatabase) SubscribeMqttTopics() {
	callback := func(client mqtt.Client, message mqtt.Message) {
		log.Debugf("localmqtt-broker topic: %s, payload: %s", message.Topic(), message.Payload())
		if len(message.Payload()) > 0 {
			covTopic := localmqtt.MakeTopic([]string{localmqtt.PointValueTopic, localmqtt.CovTopic, localmqtt.AllTopic})
			if strings.Contains(message.Topic(), covTopic) {
				d.checkAndClearPointCov(message)
			}
		}
	}
	localmqtt.SubscribeMqttTopics(callback)
}

func (d *GormDatabase) checkAndClearPointCov(message mqtt.Message) {
	topics := strings.Split(message.Topic(), "/")
	if len(topics) < 2 {
		return
	}
	pointUUID := topics[len(topics)-2]
	point, _ := d.GetPoint(pointUUID, argspkg.Args{})
	if point == nil {
		topic := message.Topic()
		log.Warnf("no point with topic: %s", topic)
		log.Warnf("clearing topic: %s, having payload: %s", topic, message.Payload())
		localmqtt.PublishValue(topic, "")
	}
}
