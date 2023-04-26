package localmqtt

import (
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/utils/boolean"
	log "github.com/sirupsen/logrus"
)

var localMqtt *LocalMqtt

func GetLocalMqtt() *LocalMqtt {
	return localMqtt
}

func Init(ip string, conf *config.Configuration, onConnected interface{}) error {
	pm := new(LocalMqtt)
	pm.QOS = mqttclient.QOS(conf.MQTT.QOS)
	pm.Retain = boolean.IsTrue(conf.MQTT.Retain)
	pm.GlobalBroadcast = boolean.IsTrue(conf.MQTT.GlobalBroadcast)
	pm.PublishPointCOV = boolean.IsTrue(conf.MQTT.PublishPointCOV)
	pm.PublishPointList = boolean.IsTrue(conf.MQTT.PublishPointList)
	pm.PointWriteListener = boolean.IsTrue(conf.MQTT.PointWriteListener)
	pm.PublishScheduleCOV = boolean.IsTrue(conf.MQTT.PublishScheduleCOV)
	pm.PublishScheduleList = boolean.IsTrue(conf.MQTT.PublishScheduleList)
	pm.ScheduleWriteListener = boolean.IsTrue(conf.MQTT.ScheduleWriteListener)
	c, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers:       []string{ip},
		Username:      conf.MQTT.Username,
		Password:      conf.MQTT.Password,
		AutoReconnect: conf.MQTT.AutoReconnect,
		ConnectRetry:  conf.MQTT.ConnectRetry,
	}, onConnected)
	if err != nil {
		log.Info("MQTT connection error:", err)
		return err
	}
	pm.Client = c
	localMqtt = pm
	err = pm.Client.Connect()
	if err != nil {
		return err
	}
	return nil
}
