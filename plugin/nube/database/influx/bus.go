package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) BusServ() {
	handlerCreated := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("INFLUX BUS PluginsCreated isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					log.Info("INFLUX BUS PluginsCreated IsDevice", " ", dev.UUID)
					//_, err = i.addPoints(dev)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus.IsPoint(e.Topic, e)
				fmt.Println("ADD POINT ON BUS")
				if err != nil {
					return
				}
				//_, err = i.addPoint(pnt)
				if err != nil {
					return
				}
				if pnt != nil {
					log.Info("INFLUX BUS PluginsCreated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsCreated,
	}
	u, _ := utils.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerCreated)
	handlerUpdated := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("INFLUX BUS PluginsUpdated isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					//_, err = i.addPoints(dev)
					log.Info("INFLUX BUS PluginsUpdated IsDevice", " ", dev.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				if pnt != nil {
					//_, err = i.pointPatch(pnt)
					log.Info("INFLUX BUS PluginsUpdated IsPoint", " ", pnt.UUID)
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsUpdated,
	}
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerUpdated)
	handlerDeleted := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				log.Info("INFLUX BUS DELETED NEW MSG", " ", e.Topic)
				//try and match is network
				net, err := eventbus.IsNetwork(e.Topic, e)
				if err != nil {
					return
				}
				if net != nil {
					log.Info("INFLUX BUS DELETED isNetwork", " ", net.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is device
				dev, err := eventbus.IsDevice(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					//_, err = i.addPoints(dev)
					log.Info("INFLUX BUS DELETED IsDevice", " ", dev.UUID)
					if err != nil {
						return
					}
					return
				}
				//try and match is point
				pnt, err := eventbus.IsPoint(e.Topic, e)
				if err != nil {
					return
				}
				log.Info("INFLUX BUS DELETED IsPoint", " ")
				if pnt != nil {
					//p, err := i.deletePoint(pnt)
					log.Info("INFLUX BUS DELETED IsPoint", " ", pnt.UUID, "WAS DELETED", " ", "p")
					if err != nil {
						return
					}
					return
				}
			}()
		},
		Matcher: eventbus.PluginsDeleted,
	}
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerDeleted)
	handlerJobs := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//p, _ := e.Data.(model.Job)
				fmt.Println("JOBS TRIGGER ON EVENT BUS")

			}()
		},
		Matcher: eventbus.JobTrigger,
	}
	u, _ = utils.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerJobs)

}
