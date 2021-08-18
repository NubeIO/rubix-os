package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/mqtt_client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)


var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MSG recieved pointsValue: %s\n", msg.Payload())
}


func main() {

	//NewAgent
	a := mqtt_client.NewClient(mqtt_client.ClientOptions{
		Servers: []string{"tcp://192.168.15.104:1883"},
	})
	err := a.Connect()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(a.IsConnected())

	err = a.Publish("adsdas", mqtt_client.AtMostOnce, false, "ddd")
	if err != nil {
		fmt.Println(err)
	}

	err = a.Subscribe("e", mqtt_client.AtMostOnce, handle)
	if err != nil {
		fmt.Println(err)
	}

	err = a.Subscribe("ee", mqtt_client.AtMostOnce, handle)
	if err != nil {
		fmt.Println(err)
	}
	err = a.Unsubscribe("e")
	if err != nil {
		fmt.Println(err)
	}

	//
	for {
		<-time.After(1 * time.Second)
	}

}
