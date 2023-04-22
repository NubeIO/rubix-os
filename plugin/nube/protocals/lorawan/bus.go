package main

import (
	"context"
	"fmt"
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
				if !inst.csConnected {
					return
				}
				p, _ := e.Data.(mqtt.Message)
				if !checkMqttTopicUplink(p.Topic()) {
					return
				}
				inst.handleMqttUplink(p)
			}()
		},
		Matcher: eventbus.MQTTUpdated,
	}
	u, _ := nuuid.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerMQTT)
}

// checkMqttTopicUplink checks the topic is a CS uplink event
func checkMqttTopicUplink(topic string) bool {
	s := strings.Split(topic, "/")
	return len(s) == 6 && s[0] == "application" && s[2] == "device" && s[4] == "event" && s[5] == "up"
}
