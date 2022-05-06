package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) BusServ() {
	handlerJobs := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				dev, err := eventbus.IsJob(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					log.Info("SYSTEM BUS PluginsCreated IsJob", " ", dev.UUID)
					if err != nil {
						return
					}
					return
				}
				return
			}()
		},
		Matcher: eventbus.JobTrigger,
	}
	u, _ := nuuid.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerJobs)

	handlerJobs = bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				dev, err := eventbus.IsSchedule(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					log.Info("SYSTEM BUS PluginsCreated IsJob", " ", dev.UUID)
					if err != nil {
						return
					}
					return
				}
				return
			}()
		},
		Matcher: eventbus.JobTrigger,
	}
	u, _ = nuuid.MakeUUID()
	key = fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerJobs)
}
