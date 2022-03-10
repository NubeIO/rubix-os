package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/mustafaturan/bus/v3"
	"strings"
)

func (i *Instance) BusServ() {
	handlerJobs := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				if strings.Split(e.Topic, ".")[2] == path {
					_, err := i.syncInflux()
					if err != nil {
						return
					}
				}
				return
			}()
		},
		Matcher: eventbus.JobTrigger,
	}
	u, _ := utils.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerJobs)
}
