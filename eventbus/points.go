package eventbus

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	"github.com/mustafaturan/bus/v3"
	"github.com/sirupsen/logrus"
)

func (eb *notificationService) registerPointsProducer() {
	handler := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			go func() {
				switch e.Topic {
				case PointCreated:
					payload, ok := e.Data.(*model.Point)
					msg := fmt.Sprintf("point %s created", payload.Name)
					logrus.Info(msg)
					if !ok {
						return
					}

				case PointUpdated:
					payload, ok := e.Data.(*model.Point)
					msg := fmt.Sprintf("event %s wiii", payload.Name)
					//publishMQTT(payload)
					logrus.Info(msg)
					if !ok {
						return
					}
				case PointCOV:
					payload, ok := e.Data.(*model.Point)
					//publishMQTT(payload)
					logrus.Error(payload)
					if !ok {
						return
					}

				}
			}()
		},
		Matcher: PointsAll,
	}
	eb.eb.RegisterHandler("points.*", handler)
}
