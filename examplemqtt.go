package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"time"
)

var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MSG recieved pointsValue: %s\n", msg.Payload())
}

func main() {

	//internalMQTT, err := mqttclient.InternalMQTT()
	//if err != nil {
	//	return
	//}
	//
	//internalMQTT.Connect()
	//internalMQTT.Publish("adsdas", mqttclient.AtMostOnce, false, "ddd")
	//NewAgent
	//a, _ := mqttclient.NewClient(mqttclient.ClientOptions{
	//	Servers: []string{"tcp://192.168.15.100:1883", "tcp://192.168.15.104:1883"},
	//})
	//err := a.Connect()
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//fmt.Println(a.IsConnected())
	//
	//err = a.Publish("adsdas", mqttclient.AtMostOnce, false, "ddd")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//err = a.Subscribe("e", mqttclient.AtMostOnce, handle)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//err = a.Subscribe("ee", mqttclient.AtMostOnce, handle)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err = a.Unsubscribe("e")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//
	for {
		<-time.After(1 * time.Second)
	}

}
