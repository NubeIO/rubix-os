package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pprint "github.com/NubeIO/lib-networking/print"
	"strconv"
	"strings"

	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
)

func (inst *Instance) BusServ() {
	handlerMQTT := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				message, _ := e.Data.(mqtt.Message)
				if messageWhois(message.Topic()) {

				}
				if messageRead(message.Topic()) {
					readType, _, _, _ := getReadType(message.Topic())
					if readType == mqttTypeReadPV {
						pv, err := decodePointPV(message)
						fmt.Println(err)
						pprint.PrintJOSN(pv)
					}

				}

			}()
		},
		Matcher: eventbus.BACnetMQTTMessage,
	}
	u, _ := nuuid.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerMQTT)
}

const (
	mqttTypeName   = "name"
	mqttTypePri    = "pri"
	mqttTypeReadPV = "pv"
)

func decode(topicType string) {
	switch topicType {
	case mqttTypeName:
		fmt.Println("one")
	case mqttTypePri:
		fmt.Println("two")
	case mqttTypeReadPV:
		fmt.Println("three")
	}
}

type payloadRawPri struct {
	Value          []string `json:"value"`
	DeviceInstance string   `json:"deviceInstance"`
}

type payloadRawRead struct {
	DeviceInstance string `json:"deviceInstance"`
	Value          string `json:"value"`
}

type payloadReadPV struct {
	IoType         string  `json:"ioType"`
	IoNumber       int     `json:"ioNumber"`
	DeviceInstance int     `json:"deviceInstance"`
	Value          float64 `json:"value"`
}

type payloadPointName struct {
	ioType         string
	ioNumber       int
	deviceInstance int
	value          string
}

func decodePointPV(msg mqtt.Message) (*payloadReadPV, error) {
	var payload *payloadRawRead
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		return nil, err
	}
	_, ioType, ioNumber, err := getReadType(msg.Topic())
	if err != nil {
		return nil, err
	}
	deviceInstance, err := s2i(payload.DeviceInstance)
	value, err := s2f(payload.Value)
	return &payloadReadPV{
		IoType:         ioType,
		IoNumber:       ioNumber,
		DeviceInstance: deviceInstance,
		Value:          value,
	}, err
}

func decodePointName(msg mqtt.Message) (*payloadPointName, error) {
	payload := new(payloadRawRead)
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		return nil, err
	}
	_, ioType, ioNumber, err := getReadType(msg.Topic())
	if err != nil {
		return nil, err
	}
	deviceInstance, err := s2i(payload.DeviceInstance)
	return &payloadPointName{
		ioType:         ioType,
		ioNumber:       ioNumber,
		deviceInstance: deviceInstance,
		value:          payload.Value,
	}, err
}

func getReadType(topic string) (topicType, ioType string, ioNumber int, err error) {
	parts := messageGetReadParts(topic)
	if len(parts) == 6 {
		v, err := strconv.Atoi(parts[4])
		if err != nil {
			return "", "", 0, errors.New("failed to convert IONumber")
		}
		return parts[5], parts[3], v, nil
	}
	return "", "", 0, errors.New("failed to decode topic")
}

func messageGetReadParts(topic string) []string {
	if messageRead(topic) {
		return strings.Split(topic, "/")
	}
	return nil
}

func messageRead(topic string) bool {
	s := strings.Split(topic, "/")
	return len(s) >= 3 && s[0] == "bacnet" && s[1] == "cmd_result" && s[2] == "read_value"
}

func messageWhois(topic string) bool {
	s := strings.Split(topic, "/")
	return len(s) == 3 && s[0] == "bacnet" && s[1] == "cmd_result" && s[2] == "whois"
}
