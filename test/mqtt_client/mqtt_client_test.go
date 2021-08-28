package mqtt_client

import (
	"fmt"
	"github.com/NubeDev/flow-framework/mqttClient"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"testing"
	"time"
)

var handle mqttClient.MessageHandler = func(client mqttClient.Client, msg mqttClient.Message) {
	fmt.Printf("MSG recieved pointsValue: %s\n", msg.Payload())
}

func TestMQTTClient(t *testing.T) {

	a := mqttClient.NewClient(mqttClient.ClientOptions{
		Servers: []string{"tcp://192.168.15.100:1883", "tcp://192.168.15.104:1883"},
	})
	err := a.Connect()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(a.IsConnected())

	err = a.Publish("adsdas", mqttClient.AtMostOnce, false, "ddd")
	if err != nil {
		fmt.Println(err)
	}

	err = a.Subscribe("e", mqttClient.AtMostOnce, handle)
	if err != nil {
		fmt.Println(err)
	}

	err = a.Subscribe("ee", mqttClient.AtMostOnce, handle)
	if err != nil {
		fmt.Println(err)
	}
	err = a.Unsubscribe("e")
	if err != nil {
		fmt.Println(err)
	}

	for {
		<-time.After(1 * time.Second)
	}
}
