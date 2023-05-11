package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/csrest"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var avoidErrorsOnStartActive = true

// checkMqttTopicCS checks the topic is a CS event
func checkMqttTopicCS(topic string) bool {
	s := strings.Split(topic, "/")
	return len(s) == 6 && s[0] == "application" && s[2] == "device" && s[4] == "event"
}

// handleMqttUplink handle CS MQTT event
func (inst *Instance) handleMqttEvent(body mqtt.Message) {
	eventType := strings.Split(body.Topic(), "/")[5]
	switch eventType {
	case "up":
		inst.handleMqttUplink(body)
	case "error":
		inst.handleMqttError(body)
	}
}

// handleMqttUplink handle CS MQTT event uplink
func (inst *Instance) handleMqttUplink(body mqtt.Message) {
	payload := new(csrest.MQTTUplink)
	json.Unmarshal(body.Payload(), &payload)
	log.Trace("lorawan: uplink: ", *payload)

	currDev, err := inst.checkAndAddValidCSDeviceFromEvent(payload.DevEUI, true)
	if err != nil {
		return
	}

	ok := inst.parseUplinkData(&payload.Object, currDev)
	if !ok {
		return
	}
	inst.updateOrCreatePoint("rssi", float64(payload.RxInfo[0].Rssi), currDev)
	inst.updateOrCreatePoint("snr", float64(payload.RxInfo[0].LoRaSNR), currDev)
	currDev.InFault = false
	currDev.Message = ""
}

// handleMqttError handle CS MQTT event error
func (inst *Instance) handleMqttError(body mqtt.Message) {
	if avoidErrorsOnStartActive {
		return
	}
	payload := new(csrest.MQTTError)
	json.Unmarshal(body.Payload(), &payload)
	currDev, err := inst.checkAndAddValidCSDeviceFromEvent(payload.DevEUI, false)
	if err != nil {
		return
	}
	currDev.InFault = true
	currDev.Message = fmt.Sprintf("%s: %s", payload.Type, payload.Error)
	inst.db.UpdateDeviceErrors(currDev.UUID, currDev)
}

func (inst *Instance) checkAndAddValidCSDeviceFromEvent(devEUI string, withPoints bool) (*model.Device, error) {
	currDev, err := inst.db.GetDeviceByArgs(api.Args{AddressUUID: &devEUI, WithPoints: withPoints})
	if err != nil {
		var csDev *csrest.DeviceSingle
		csDev, err = inst.chirpStack.GetDevice(devEUI)
		if csrest.IsCSConnectionError(err) {
			inst.setCSDisconnected(err)
			return nil, err
		}
		if err != nil {
			log.Debug("lorawan: MQTT event recived but device missing from Chirpstack. EUI=%s, Error: %s", devEUI, err)
			return nil, err
		}
		log.Debug("lorawan: adding new device from uplink")
		currDev, err = inst.addMissingDeviceSingle(csDev)
	}
	return currDev, err
}

func (inst *Instance) parseUplinkData(data *map[string]interface{}, device *model.Device) bool {
	log.Trace(fmt.Sprintf("lorawan: parsing uplink for device UUID=%s, EUI=%s, name=%s", device.UUID, *device.AddressUUID, device.Name))
	if len(*data) == 0 {
		return false
	}
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
		case string:
			str := reflect.ValueOf(v).String()
			value, err = strconv.ParseFloat(str, 64)
			if err != nil {
				value, err = nstring.ConvertKnownStringToFloat(str)
			}
			if err != nil {
				log.Debug("lorawan: could not parse string to float: %s", reflect.ValueOf(v).String())
				continue
			}

		default:
			log.Debug(fmt.Sprintf("lorawan: parseUplinkData unsupported value type: %T = %v", t, v))
			continue
		}
		inst.updateOrCreatePoint(k, value, device)
	}
	return true
}

func (inst *Instance) updateOrCreatePoint(pointName string, value float64, device *model.Device) (err error) {
	point := inst.getPointByAddressUUID(pointName, *device.AddressUUID, device.Points)
	if point == nil {
		point, err = inst.createNewPoint(pointName, *device.AddressUUID, device.UUID)
		if err != nil {
			return err
		}
	}
	log.Trace(fmt.Sprintf("lorawan: update point %s value=%f", *point.AddressUUID, value))
	inst.pointWrite(point.UUID, value)
	return nil
}
