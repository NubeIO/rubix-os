package eventbus

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/mqtt_client"
	"log"
)

func publishMQTT(sensorStruct *model.Point) {
	a := mqtt_client.NewClient(mqtt_client.ClientOptions{
		Servers: []string{"tcp://0.0.0.0:1883"},
	})
	err := a.Connect();if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a.IsConnected())
	topic := fmt.Sprintf("rubix/%s", sensorStruct.UUID)
	data, err := json.Marshal(sensorStruct);if err != nil {
		log.Fatal(err)
	}

	err = a.Publish(topic, mqtt_client.AtMostOnce, false, string(data));if err != nil {
		log.Fatal(err)
	}

}


