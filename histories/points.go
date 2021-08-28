package histories

import (
	"math"
)

//PointCOV event for the point value
func PointCOV(new float64, existingData float64, cov float64) (bool, float64) {
	c := new - existingData
	if math.Abs(c) >= cov {
		return false, existingData
	} else {
		return false, new
	}
}


//func PublishSensor(commonSensorData decoder.CommonValues, sensorStruct interface{}, mqttConn *mqtt_lib.MqttConnection) {
//	jsonValue, _ := json.Marshal(sensorStruct)
//	PublishJSON(commonSensorData, jsonValue, mqttConn)
//}
//
//func PublishJSON(commonSensorData decoder.CommonValues, jsonValue []byte, mqttConn *mqtt_lib.MqttConnection) {
//	topic := "test-topic/" + string(commonSensorData.Id)
//	log.Printf("MQTT PUB: {\"topic\": \"%s\", \"payload\": \"%s\"}", topic, string(jsonValue))
//	mqttConn.Publish(string(jsonValue), topic)
//}