package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NubeIO/flow-framework/services/localmqtt"
	"github.com/NubeIO/flow-framework/src/gocancel"
	"github.com/NubeIO/flow-framework/utils/float"
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
	inst.rubixpointsyncDebugMsg("subscribeToMQTTForPointCOV()")
	callback := func(client mqtt.Client, message mqtt.Message) {
		body := &covPayload{}
		err := json.Unmarshal(message.Payload(), &body)
		if err == nil {
			if body != nil {
				// d.publishPointWrite(body)
				messageTopic := message.Topic()
				requiredNetworksArray := inst.config.Job.Networks
				if requiredNetworksArray == nil || len(requiredNetworksArray) == 0 {
					requiredNetworksArray = []string{"system"}
				}
				for _, reqNetwork := range requiredNetworksArray {
					topicParts := strings.Split(messageTopic, "/")
					netName := topicParts[7]
					devName := topicParts[9]
					pntName := topicParts[11]
					if netName == reqNetwork { // topicParts[7] is the network name
						inst.rubixpointsyncDebugMsg(fmt.Sprintf("subscribeToMQTTForPointCOV() message: %+v", message))
						inst.rubixpointsyncDebugMsg(fmt.Sprintf("subscribeToMQTTForPointCOV() body: %+v", body))
						pointValue := float.New(body.Value)
						err := inst.SyncSingleRubixPointWithFF(netName, devName, pntName, pointValue) // topicParts[10] is the point UUID
						if err != nil {
							inst.rubixpointsyncDebugMsg("subscribeToMQTTForPointCOV() error:", err)
						}
					}
				}
			}
		}
	}
	// topic := fetchPointsTopicWrite
	// rubix/points/value/cov/all/#
	var topic = "rubix/points/value/cov/all/#"
	mqttClient := localmqtt.GetLocalMqtt().Client
	if mqttClient != nil {
		err := mqttClient.Subscribe(topic, 1, callback)
		if err != nil {
			inst.rubixpointsyncErrorMsg(fmt.Sprintf("localmqtt-broker subscribe:%s err:%s", topic, err.Error()))
		} else {
			inst.rubixpointsyncDebugMsg(fmt.Sprintf("localmqtt-broker subscribe:%s", topic))
		}
	}
}

func (inst *Instance) StartMQTTSubscribeCOV() error {
	ctx, cancel := context.WithCancel(context.Background())
	inst.mqttCancel = cancel
	go gocancel.GoRoutineWithContextCancel(ctx, inst.subscribeToMQTTForPointCOV)
	return nil
}
