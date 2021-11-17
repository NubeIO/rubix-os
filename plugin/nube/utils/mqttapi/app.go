package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/mqttclient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

//connect to brokers
func (i *Instance) connect() {
	log.Info("mqtt-api: sync has is been called")

	integrations, err := i.db.GetEnabledIntegrationByPluginConfId(i.pluginUUID)
	if err != nil {
		//return false, err
	}
	for _, integration := range integrations {
		fmt.Println(integration)

	}
	if len(integrations) == 0 {
		log.Info("mqtt-api: can't be registered, integration details missing.")
	}

	localBroker := "tcp://0.0.0.0:1883" //TODO add to config, this is meant to be an unsecure broker
	connected, err := i.initMQTT(localBroker)
	if err != nil {
		log.Error(err, "mqtt-api: failed to broker")
	}
	fmt.Println(connected, "connected")
	i.subscribe()

}

//used for getting data into the plugins
var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Println(msg.Topic(), " ", "NEW MQTT MES", " ", string(msg.Payload()))

}

func (i *Instance) subscribe() {
	err := i.mqtt.Subscribe("test", mqttclient.AtMostOnce, handle) //lorawan chirpstack
	if err != nil {

	}

}

//initMQTT  mqtt connection
func (i *Instance) initMQTT(ip string) (bool, error) {
	c, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{ip},
	})
	if err != nil {
		log.Error(err, "mqtt-api: failed to broker")
	}
	log.Info(err, "CONNECT to broker")
	i.mqtt = c

	err = c.Connect()
	if err != nil {
		log.Error(err, "mqtt-api: failed to broker")
	} else {
		fmt.Println(err)
	}
	return c.IsConnected(), nil
}
