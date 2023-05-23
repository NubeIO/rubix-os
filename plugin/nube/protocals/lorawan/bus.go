package main

import (
	"context"
	"fmt"
	"time"

	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
)

var csBusKey string = ""

func (inst *Instance) busServ() {
	handlerMQTT := bus.Handler{
		Handle: func(_ context.Context, e bus.Event) {
			if !inst.csConnected {
				return
			}
			p, _ := e.Data.(mqtt.Message)
			if !checkMqttTopicCS(p.Topic()) {
				return
			}
			inst.handleMqttEvent(p)
		},
		Matcher: eventbus.MQTTUpdated,
	}
	u, _ := nuuid.MakeUUID()
	csBusKey = fmt.Sprintf("key_%s", u)
	avoidErrorsOnStartActive = true
	eventbus.GetBus().RegisterHandler(csBusKey, handlerMQTT)
	go func() {
		time.Sleep(5 * time.Second)
		avoidErrorsOnStartActive = false
	}()
}

func (inst *Instance) busDisable() {
	eventbus.GetBus().DeregisterHandler(csBusKey)
}
