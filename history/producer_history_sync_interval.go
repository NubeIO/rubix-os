package history

import (
	"encoding/json"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

func (h *History) InitProducerHistorySyncInterval(syncPeriod int) {
	h.cron = gocron.NewScheduler(time.UTC)
	_, _ = h.cron.Every(syncPeriod).Tag("ProducerHistorySyncInterval").Do(h.syncProducerHistoryInterval)
	h.cron.StartAsync()
}

func (h *History) syncProducerHistoryInterval() {
	log.Info("Producer history sync interval has is been called...")
	producers, err := h.DB.GetProducersForCreateInterval()
	if err != nil {
		log.Error(err)
		return
	}
	currentDate := time.Now().UTC()
	for _, producer := range producers {
		if producer.HistoryInterval == nil || *producer.HistoryInterval < 1 {
			continue
		}
		if producer.ProducerThingClass != "point" { // TODO: CreateProducerHistory for ProducerThingClass == "schedule"
			continue
		}
		latestPH, _ := h.DB.GetLatestProducerHistoryByProducerUUID(producer.UUID)
		if latestPH == nil || currentDate.Sub(latestPH.Timestamp).Seconds() >= float64(*producer.HistoryInterval*60) {
			latestPH = new(model.ProducerHistory)
			// Minutes is placing such a way if 15, then it will store values on 0, 15, 30, 45
			minute := (currentDate.Minute() / *producer.HistoryInterval) * *producer.HistoryInterval
			latestPH.ProducerUUID = producer.UUID
			latestPH.Timestamp = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				currentDate.Hour(), minute, 0, 0, time.UTC)
			point, _ := h.DB.GetPoint(producer.ProducerThingUUID, api.Args{WithPriority: true})
			if point != nil {
				latestPH.DataStore, _ = json.Marshal(point.Priority)
				_, err := h.DB.CreateProducerHistory(latestPH)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
	log.Info("Finished producer history sync interval process")
}
