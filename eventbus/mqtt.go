package eventbus

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/mqttclient"
	"log"
)

func publishMQTT(sensorStruct model.ProducerBody) {
	a := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err := a.Connect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a.IsConnected())
	topic := fmt.Sprintf("rubix/%s", sensorStruct.ProducerUUID)
	fmt.Println(11111, topic)
	data, err := json.Marshal(sensorStruct)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Publish(topic, mqttclient.AtMostOnce, false, string(data))
	if err != nil {
		log.Fatal(err)
	}

}
