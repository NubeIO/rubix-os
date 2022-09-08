package main

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csmodel"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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
	log.Debug("lorawan: Uplink: ", *payload)

	currDev, err := inst.db.GetDeviceByArgs(api.Args{AddressUUID: &payload.DevEUI, WithPoints: true})
	if err != nil || currDev == nil {
		csDev, err := inst.REST.GetDevice(payload.DevEUI)
		if csrest.IsCSConnectionError(err) {
			inst.setCSDisconnected(err)
			return
		}
		if err != nil {
			log.Warnf("lorawan: MQTT Uplink recived but device missing from Chirpstack. EUI=%s, Error: %s", payload.DevEUI, err)
			return
		}
		log.Debug("lorawan: Adding new device from uplink")
		currDev, err = inst.createDeviceFromCSDevice(csDev)
		if err != nil {
			return
		}
	}
	inst.parseUplinkData(&payload.Object, currDev)
}

// checkMqttTopicUplink checks the topic is a CS uplink event
func checkMqttTopicUplink(topic string) bool {
	s := strings.Split(topic, "/")
	return len(s) == 6 && s[0] == "application" && s[2] == "device" && s[4] == "event" && s[5] == "up"
}

func (inst *Instance) parseUplinkData(data *map[string]interface{}, device *model.Device) {
	log.Debugf("lorawan: Parsing uplink for device UUID=%s, EUI=%s, name=%s", device.UUID, *device.AddressUUID, device.Name)
	var err error = nil
	for k, v := range *data {
		var value float64
		switch t := v.(type) {
		case int:
			value = float64(reflect.ValueOf(v).Int())
		case float64:
			value = float64(reflect.ValueOf(v).Float())
		case float32:
			value = float64(reflect.ValueOf(v).Float())
		case bool:
			if reflect.ValueOf(v).Bool() {
				value = 1
			} else {
				value = 0
			}
		case map[string]interface{}:
			dataInternal := v.(map[string]interface{})
			inst.parseUplinkData(&dataInternal, device)

		default:
			log.Warnf("lorawan: parseUplinkData unsupported value type: %T = %t", t, v)
			continue
		}
		point := inst.getPointByAddressUUID(k, *device.AddressUUID, device.Points)
		if point == nil {
			point, err = inst.createNewPoint(k, *device.AddressUUID, device.UUID)
			if err != nil {
				continue
			}
		}
		log.Debugf("lorawan: Update point %s value=%f", *point.AddressUUID, value)
		inst.pointWrite(point.UUID, value)
	}
}
