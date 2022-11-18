package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-framework/services/localmqtt"
	"github.com/NubeIO/flow-framework/src/gocancel"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

type covPayload struct {
	Value    float64 `json:"value"`
	ValueRaw float64 `json:"value_raw"`
	Ts       string  `json:"ts"`
	Priority int     `json:"priority"`
}

func (inst *Instance) subscribeToMQTTForPointCOV() {
	inst.edgeazureDebugMsg("subscribeToMQTTForPointCOV()")
	callback := func(client mqtt.Client, message mqtt.Message) {
		body := &covPayload{}
		err := json.Unmarshal(message.Payload(), &body)
		if err == nil {
			if body != nil {
				// d.publishPointWrite(body)
				messageTopic := message.Topic()
				inst.edgeazureDebugMsg("subscribeToMQTTForPointCOV() messageTopic:", messageTopic)
				pluginsArray := inst.config.Job.Networks
				if pluginsArray == nil || len(pluginsArray) == 0 {
					pluginsArray = []string{"system"}
				}
				for _, plugin := range pluginsArray {
					topicParts := strings.Split(messageTopic, "/")
					if topicParts[5] == plugin { // topicParts[5] is the plugin name
						inst.edgeazureDebugMsg(fmt.Sprintf("subscribeToMQTTForPointCOV() message: %+v", message))
						inst.edgeazureDebugMsg(fmt.Sprintf("subscribeToMQTTForPointCOV() body: %+v", body))
						inst.edgeazureDebugMsg("subscribeToMQTTForPointCOV() topicParts[10]:", topicParts[10])
						/*
							err := inst.SendPointWriteHistory(topicParts[10]) // topicParts[10] is the point UUID
							if err != nil {
								inst.edgeazureErrorMsg("subscribeToMQTTForPointCOV() error:", err)
							}
						*/
					}
				}
			}
		}
	}
	// topic := fetchPointsTopicWrite
	// rubix/points/value/cov/all/#
	var topic = "rubix/points/value/cov/all/#"
	mqttClient := localmqtt.GetPointMqtt().Client
	if mqttClient != nil {
		err := mqttClient.Subscribe(topic, 1, callback)
		if err != nil {
			inst.edgeazureErrorMsg(fmt.Sprintf("localmqtt-broker subscribe:%s err:%s", topic, err.Error()))
		} else {
			inst.edgeazureDebugMsg(fmt.Sprintf("localmqtt-broker subscribe:%s", topic))
		}
	}
}

func (inst *Instance) StartMQTTSubscribeCOV() error {
	ctx, cancel := context.WithCancel(context.Background())
	inst.mqttCancel = cancel
	go gocancel.GoRoutineWithContextCancel(ctx, inst.subscribeToMQTTForPointCOV)
	return nil
}
