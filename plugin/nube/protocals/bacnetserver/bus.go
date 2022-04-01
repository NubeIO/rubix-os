package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) BusServ() {

	handlerUpdated := bus.Handler{ //UPDATED
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is point
				pnt, err := eventbus.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				if pnt != nil {
					_, err = inst.updatePointValue(pnt)
					log.Info("BACNET-SERVER BUS PluginsUpdated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsUpdated,
	}
	u, _ := utils.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerUpdated)

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
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerMQTT)

}
