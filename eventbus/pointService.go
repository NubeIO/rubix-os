package eventbus

import (
	"context"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/mustafaturan/bus/v3"
	"github.com/sirupsen/logrus"
)

func (eb *notificationService) registerPointsSubscriber() {
	handler := bus.Handler{
		Handle:  func(ctx context.Context, e bus.Event){
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
					msg := fmt.Sprintf("payment %s paid", payload.Name)
					logrus.Info(msg)

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





