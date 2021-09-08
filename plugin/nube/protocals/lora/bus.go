package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (c *Instance) BusServ()  {
	pluginUUID := c.pluginUUID
	pluginHandler := bus.Handler{
		Handle:  func(ctx context.Context, e bus.Event){
			go func() {
				switch e.Topic {
				case pluginUUID:
					fmt.Println("plugin event")
					payload := e.Data
					msg := fmt.Sprintf("point %s created", payload)
					log.Info(msg)
				}
			}()
		},
		Matcher: eventbus.PluginsAll,
	}
	keyP := fmt.Sprintf("key_%s", pluginUUID)
	eventbus.GetBus().RegisterHandler(keyP, pluginHandler)

	networkUUID := c.networkUUID
	handler := bus.Handler{
		Handle:  func(ctx context.Context, e bus.Event){
			go func() {
				switch e.Topic {
				case networkUUID:
					fmt.Println("network event")
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
		Matcher: eventbus.NetworksAll,
	}
	keyN := fmt.Sprintf("key_%s", pluginUUID)
	eventbus.GetBus().RegisterHandler(keyN, handler)

}
