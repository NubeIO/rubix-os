package main

import (
	"encoding/json"
	"strings"

	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// handleMqttUplink parse CS MQTT uplink data
func (inst *Instance) handleMqttUplink(body mqtt.Message) {
	payload := new(csmodel.BaseUplink)
	err := json.Unmarshal(body.Payload(), &payload)
	if err != nil {
		log.Error("lorawan: Invalid MQTT uplink data: ", err)
		return
	}
	log.Debug(*payload)

	if !inst.euiExists(payload.DevEUI) {
		dev, err := inst.REST.GetDevice(payload.DevEUI)
		if csrest.IsCSConnectionError(err) {
			inst.setCSDisconnected(err)
			return
		}
		if err != nil {
			log.Warn("lorawan: MQTT Uplink recived but device missing from Chirpstack ", payload.DevEUI, " ", err)
			return
		}
		log.Info("lorawan: NEW DEVICE TO ADD: ", *payload, *dev)
	}
	// TODO: update device points
}

// checkMqttTopicUplink checks the topic is a CS uplink event
func checkMqttTopicUplink(topic string) bool {
	s := strings.Split(topic, "/")
	if len(s) != 6 ||
		!(s[0] == "application" && s[2] == "device" && s[4] == "event" && s[5] == "up") {
		return false
	}
	return true
}
