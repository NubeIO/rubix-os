package history

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

func (h *History) InitIntervalHistoryCreator(syncPeriod int) {
	h.cron = gocron.NewScheduler(time.UTC)
	_, _ = h.cron.Every(syncPeriod).Tag("InitIntervalHistoryCreator").Do(h.createIntervalHistories)
	h.cron.StartAsync()
}

func (h *History) createIntervalHistories() {
	log.Info("Create interval histories has been called...")
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
			// minutes are placing such a way if 15, then it will store values on 0, 15, 30, 45
			minute := (currentDate.Minute() / *producer.HistoryInterval) * *producer.HistoryInterval
			latestPH.ProducerUUID = producer.UUID
			latestPH.Timestamp = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				currentDate.Hour(), minute, 0, 0, time.UTC)
			point, _ := h.DB.GetPoint(producer.ProducerThingUUID, api.Args{})
			if point != nil {
				latestPH.PresentValue = point.PresentValue
				_, err := h.DB.CreateProducerHistory(latestPH)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
	log.Info("Finished create interval histories process")
}
