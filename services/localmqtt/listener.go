package localmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func SubscribeMqttTopics(handler mqtt.MessageHandler) {
	qos := byte(localMqtt.QOS)
	filters := map[string]byte{
		MakeTopic([]string{PointValueTopic, CovTopic, AllTopic, MultiLevelWildcard}): qos,
	}
	err := localMqtt.Client.SubscribeMultiple(filters, handler)
	if err != nil {
		log.Error(err)
		return
	} else {
		for topic, _ := range filters {
			log.Infof("localmqtt-broker subscribe: %s", topic)
		}
	}
}
