package payload

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
)

func PublishSensor(commonSensorData decoder.CommonValues, sensorStruct interface{}) {
	jsonValue, _ := json.Marshal(sensorStruct)
	fmt.Println(jsonValue)
	//PublishJSON(commonSensorData, jsonValue, mqttConn)
}

//func PublishJSON(commonSensorData decoder.CommonValues, jsonValue []byte, mqttConn *mqtt_lib.MqttConnection) {
//	topic := "test-topic/" + string(commonSensorData.Id)
//	log.Printf("MQTT PUB: {\"topic\": \"%s\", \"payload\": \"%s\"}", topic, string(jsonValue))
//	mqttConn.Publish(string(jsonValue), topic)
//}
