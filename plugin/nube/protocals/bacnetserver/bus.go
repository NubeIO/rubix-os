package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
	"strings"
)

// localmqtt topic
// bacnet/program/DEVICE-ID/state

func checkMqttTopicUplink(topic string) bool {
	if strings.Contains(topic, "bacnet/program") {
		return true

	}
	return false
}

var bacnetStarted bool

// handleMqttUplink parse bacnet-server-c
func (inst *Instance) handleMqttUplink(body mqtt.Message) {
	if string(body.Payload()) == "started" {
		log.Infof("bacnet-server just restarted topic:%s", body.Topic())
		bacnetStarted = true
	}
}

func (inst *Instance) BusServ() {
	handlerMQTT := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
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
