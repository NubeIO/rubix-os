package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (c *Instance) BusServ()  {
	topic := c.pluginUUID
	handler := bus.Handler{
		Handle:  func(ctx context.Context, e bus.Event){
			go func() {
				fmt.Println(3333, e.Topic, topic)
				switch e.Topic {
				case topic:
					payload := e.Data
					msg := fmt.Sprintf("point %s created", payload)
					log.Info(msg)
				case "**":
					//payload, ok := e.Data.(*model.Point)
					//msg := fmt.Sprintf("event %s wiii", payload.Name)
					////publishMQTT(payload)
					//logrus.Info(msg)
					//if !ok {
					//	return
					//}
				}
			}()
		},
		Matcher: topic,
	}
	key := fmt.Sprintf("key_%s", topic)
	eventbus.GetBus().RegisterHandler(key, handler)

}
