package eventbus

import (
	"github.com/NubeIO/rubix-os/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var debug = false

// used for getting data into the plugins
var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if debug {
		log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
	}
	GetService().RegisterTopic(MQTTUpdated)
	err := GetService().Emit(CTX(), MQTTUpdated, msg)
	if err != nil {
		return
	}
}

func RegisterMQTTBus(enableDebug bool) {
	if enableDebug {
		debug = true
	}
	c, _ := mqttclient.GetMQTT()
	// TODO this needs to be removed as its for a plugin, the plugin needs to register the topics it wants the main framework to subscribe to, also unsubscribe when the plugin is disabled
	c.Subscribe("+/+/+/+/+/+/rubix/bacnet_server/points/+/#", mqttclient.AtMostOnce, handle) // bacnet-server
	c.Subscribe("+/+/+/+/+/+/rubix/bacnet_master/points/+/#", mqttclient.AtMostOnce, handle) // bacnet-bserver
	c.Subscribe("application/+/device/+/event/+", mqttclient.AtMostOnce, handle)             // lorawan
	c.Subscribe("bacnet/program/#", mqttclient.AtMostOnce, handle)                           // bacnet-server
}
