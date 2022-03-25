package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
)

func (inst *Instance) BusServ() {
	handlerMQTT := bus.Handler{ //MQTT UPDATE (got as msg over from bacnet stack)
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				p, _ := e.Data.(mqtt.Message)
				_, err := inst.bacnetUpdate(p)
				if err != nil {
					return
				}
			}()
		},
		Matcher: eventbus.MQTTUpdated,
	}
	u, _ := utils.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerMQTT)

}
