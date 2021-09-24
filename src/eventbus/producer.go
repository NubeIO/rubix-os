package eventbus

import (
	"context"
	"github.com/NubeDev/flow-framework/model"
	"github.com/mustafaturan/bus/v3"
)

func (eb *notificationService) registerProducer() {
	handler := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				switch e.Topic {
				case ProducerEvent:
					payload, ok := e.Data.(model.ProducerBody)
					publishMQTT(payload)
					if !ok {
						return
					}
				}
			}()
		},
		Matcher: ProducerAll,
	}
	eb.eb.RegisterHandler(ProducerAll, handler)
}
