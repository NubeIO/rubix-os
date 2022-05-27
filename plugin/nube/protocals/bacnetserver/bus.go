package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) BusServ() {

	handlerUpdated := bus.Handler{ // UPDATED
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {

				isThis := eventbus.IsThisPlugin(e.Topic, inst.pluginUUID)
				if !isThis {
					return
				}
				// try and match is point
				pnt, err := eventbus.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				if pnt != nil {
					log.Info("BACNET-SERVER BUS PluginsUpdated IsPoint", " ", pnt.UUID)
					_, err = inst.updatePointValue(pnt)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsUpdated,
	}
	u, _ := nuuid.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerUpdated)

	handlerMQTT := bus.Handler{ // MQTT UPDATE (got as msg over from bacnet stack)
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
	u, _ = nuuid.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerMQTT)

}
