package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/mustafaturan/bus/v3"
)

func (i *Instance) BusServ() {
	handlerJobs := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				_, err := i.syncHistory()
				if err != nil {
					return
				}
			}()
		},
		Matcher: eventbus.JobTrigger,
	}
	u, _ := utils.MakeUUID()
	key := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(key, handlerJobs)

}
