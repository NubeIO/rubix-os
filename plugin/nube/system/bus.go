package main

import (
	"context"
	"fmt"
	eventbus2 "github.com/NubeDev/flow-framework/src/eventbus"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) BusServ() {
	handlerCreated := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus2.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("LORA BUS PluginsCreated isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus2.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					log.Info("LORA BUS PluginsCreated IsDevice", " ", dev.UUID)
					//_, err = i.addPoints(dev)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus2.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				if pnt != nil {
					log.Info("LORA BUS PluginsCreated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus2.PluginsCreated,
	}
	u, _ := utils.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus2.GetBus().RegisterHandler(key, handlerCreated)
	handlerUpdated := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus2.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("LORA BUS PluginsUpdated isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus2.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					//_, err = i.addPoints(dev)
					log.Info("LORA BUS PluginsUpdated IsDevice", " ", dev.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus2.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				if pnt != nil {
					//_, err = c.addPoints(dev)
					log.Info("LORA BUS PluginsUpdated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus2.PluginsUpdated,
	}
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus2.GetBus().RegisterHandler(key, handlerUpdated)

}
