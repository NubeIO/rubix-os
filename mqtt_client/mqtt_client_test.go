package mqtt_client

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)


var handle mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MSG recieved pointsValue: %s\n", msg.Payload())
}


func main() {

	//NewAgent
	a := NewClient(ClientOptions{
		Servers: []string{"tcp://192.168.15.100:1883", "tcp://192.168.15.104:1883"},
	})
	err := a.Connect()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(a.IsConnected())

	err = a.Publish("adsdas", AtMostOnce, false, "ddd")
	if err != nil {
		fmt.Println(err)
	}

	err = a.Subscribe("e", AtMostOnce, handle)
	if err != nil {
		fmt.Println(err)
	}

	err = a.Subscribe("ee", AtMostOnce, handle)
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

