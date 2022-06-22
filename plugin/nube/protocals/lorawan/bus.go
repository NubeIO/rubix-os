package main

import (
	"context"
	"fmt"

	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
)

func (inst *Instance) BusServ() {
	handlerMQTT := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				p, _ := e.Data.(mqtt.Message)
				if !CheckMQTTTopicUplink(p.Topic()) {
					return
				}
				_, err := inst.HandleMQTTUplink(p)
				if err != nil {
					return
				}
			}()
		},
		Matcher: eventbus.MQTTUpdated,
	}
	u, _ := nuuid.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerMQTT)
}
