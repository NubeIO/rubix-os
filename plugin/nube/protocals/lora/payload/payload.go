package payload

import (
	"encoding/json"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	log "github.com/sirupsen/logrus"
)

func PublishSensor(commonSensorData decoder.CommonValues, sensorStruct interface{}) {
	jsonValue, _ := json.Marshal(sensorStruct)
	log.Info("LORA: ", string(jsonValue))
	//PublishJSON(commonSensorData, jsonValue, mqttConn)
}

//func PublishJSON(commonSensorData decoder.CommonValues, jsonValue []byte, mqttConn *mqtt_lib.MqttConnection) {
//	topic := "test-topic/" + string(commonSensorData.Id)
//	log.Printf("MQTT PUB: {\"topic\": \"%s\", \"payload\": \"%s\"}", topic, string(jsonValue))
//	mqttConn.Publish(string(jsonValue), topic)
//}
