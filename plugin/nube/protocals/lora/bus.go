package main

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/mustafaturan/bus/v3"
	"strings"
)

func getTopicPart(topic string, index int, contains string) string {
	s := strings.Split(topic, ".")
	for i, e := range s {
		if i == index {
			if strings.Contains(e, contains) { // if topic has pnt (is uuid of point)
				return e
			}
		}
	}
	return ""
}

//check if the payload is of type device
func isDev(topic string, e bus.Event) (*model.Device, error) {
	if getTopicPart(topic, 3, "dev") != "" {
		payload, _ := e.Data.(*model.Device)
		return payload, nil
	}
	return nil, nil
}

func (c *Instance) BusServ() {
	handler := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				//try and match is device
				dev, err := isDev(e.Topic, e)
				if err != nil {
					return
				}
				if dev != nil {
					_, err = c.addPoints(dev)
					if err != nil {
						return
					}
				}
				fmt.Println(e.Topic, "topic", dev)
			}()
		},
		Matcher: eventbus.PluginsAll,
	}
	u, _ := utils.MakeUUID()
	keyN := fmt.Sprintf("key_%s", u)
	eventbus.GetBus().RegisterHandler(keyN, handler)

}
