package main

import (
	"encoding/json"
	"strings"

	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	bm "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// mqttUpdate listen on mqtt and then update the point in flow-framework
func (inst *Instance) HandleMQTTUplink(body mqtt.Message) (*bm.Point, error) {

	log.Infof("lorawan: MQTT uplink")
	// log.Infof("lorawan: MQTT uplink from dev-id:%s", devEUI)

	payload := new(csmodel.BaseUplink)
	err := json.Unmarshal(body.Payload(), &payload)
	if err != nil {
		log.Errorf("lorawan: Invalid MQTT uplink data:", err)
		return nil, err
	}
	log.Print(*payload)

	// dev, err := inst.REST.GetDevice(payload.DevEUI)

	// if err != nil {
	//     log.Errorf("lorawan: check device on chirpstack exists dev-id:%s  %v", devEUI, err)
	//     return nil, err
	// }
	// // check the payload for how to decode from
	// if dev.Device.DeviceProfileID == elsysAPB {
	//     decoded := new(model.ElsysAPB)
	//     err = json.Unmarshal(body.Payload(), &decoded)
	// }
	// if err != nil {
	//     return nil, err
	// }

	return nil, nil
}

func CheckMQTTTopicUplink(topic string) bool {
	s := strings.Split(topic, "/")
	if len(s) != 6 ||
		!(s[0] == "application" && s[2] == "device" && s[4] == "event" && s[5] == "up") {
		return false
	}
	return true
}
