package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pprint "github.com/NubeIO/lib-networking/print"
	"strconv"
	"strings"
	"time"

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
					devicesFound, err := decodeWhois(message)
					fmt.Println(err)
					pprint.PrintJOSN(devicesFound)
				}
				if messageRead(message.Topic()) {
					readType, _, _, _ := getReadType(message.Topic())
					if readType == mqttTypeReadPV { // present value
						pv, err := decodePointPV(message)
						fmt.Println(err)
						pprint.PrintJOSN(pv)
						fmt.Println("SET IN STORE", pv.TxnNumber)
						if pv != nil {
							inst.store.Set(pv.TxnNumber, pv, 10*time.Second)
						}
					}
					if readType == mqttTypePri { // priority array
						fmt.Println(string(message.Payload()))
						pri, err := decodePointPri(message)
						fmt.Println(err)
						pprint.PrintJOSN(pri)
					}
					if readType == mqttTypeName { // pointName
						details, err := decodePointName(message)
						fmt.Println(err)
						pprint.PrintJOSN(details)
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

type devices struct {
	Devices []device `json:"devices"`
}

type whoIsRaw struct {
	Value []deviceRaw `json:"value"`
}

type device struct {
	DeviceId      int    `json:"device_id"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	NetworkNumber int    `json:"network_number"`
	Apdu          int    `json:"apdu"`
}

type deviceRaw struct {
	DeviceId   string `json:"device_id"`   // 1
	MacAddress string `json:"mac_address"` // 192.168.15.10:47808
	Snet       string `json:"snet"`        // network number
	Sadr       string `json:"sadr"`
	Apdu       string `json:"apdu"` // 1476, 480, 206
}

type payloadRawPri struct {
	IoType         string   `json:"ioType"`
	IoNumber       int      `json:"ioNumber"`
	Value          []string `json:"value"`
	DeviceInstance string   `json:"deviceInstance"`
}

type payloadPri struct {
	IoType         string    `json:"ioType"`
	IoNumber       int       `json:"ioNumber"`
	Value          *PriArray `json:"value"`
	DeviceInstance int       `json:"deviceInstance"`
}

type payloadRawRead struct {
	DeviceInstance string `json:"deviceInstance"`
	Value          string `json:"value"`
	TxnSource      string `json:"txn_source"`
	TxnNumber      string `json:"txn_number"`
}

type payloadReadPV struct {
	IoType         string  `json:"ioType"`
	IoNumber       int     `json:"ioNumber"`
	DeviceInstance int     `json:"deviceInstance"`
	Value          float64 `json:"value"`
	TxnSource      string  `json:"txn_source"`
	TxnNumber      string  `json:"txn_number"`
}

type payloadPointName struct {
	IoType         string `json:"ioType"`
	IoNumber       int    `json:"ioNumber"`
	DeviceInstance int    `json:"deviceInstance"`
	Value          string `json:"value"`
}

func decodeWhois(msg mqtt.Message) (*devices, error) {
	var payload *whoIsRaw
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		return nil, err
	}
	var out *devices
	for _, dev := range payload.Value {
		ip, port := decodeMac(dev.MacAddress)
		newDev := device{
			DeviceId:      s2iNoErr(dev.DeviceId),
			Ip:            ip,
			Port:          port,
			NetworkNumber: s2iNoErr(dev.Snet),
			Apdu:          s2iNoErr(dev.Apdu),
		}
		out.Devices = append(out.Devices, newDev)
	}
	return out, err
}

func decodePointPri(msg mqtt.Message) (*payloadPri, error) {
	var payload *payloadRawPri
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		return nil, err
	}
	_, ioType, ioNumber, err := getReadType(msg.Topic())
	if err != nil {
		return nil, err
	}
	deviceInstance, err := s2i(payload.DeviceInstance)

	return &payloadPri{
		IoType:         ioType,
		IoNumber:       ioNumber,
		DeviceInstance: deviceInstance,
		Value:          cleanArray(payload.Value),
	}, err
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
		TxnSource:      payload.TxnSource,
		TxnNumber:      payload.TxnNumber,
	}, err
}

func decodePointName(msg mqtt.Message) (*payloadPointName, error) {
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
	return &payloadPointName{
		IoType:         ioType,
		IoNumber:       ioNumber,
		DeviceInstance: deviceInstance,
		Value:          payload.Value,
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
