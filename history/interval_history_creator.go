package history

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

func (h *History) InitIntervalHistoryCreator(syncPeriod int) {
	h.cron = gocron.NewScheduler(time.UTC)
	h.cron.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, _ = h.cron.Every(syncPeriod).Tag("InitIntervalHistoryCreator").Do(h.createIntervalHistories)
	h.cron.StartAsync()
}

func (h *History) createIntervalHistories() {
	log.Debug("Create interval histories has been called...")
	var producerHistories []*model.ProducerHistory
	currentDate := time.Now().UTC()
	producers, err := h.DB.GetProducersForCreateInterval()
	if err != nil {
		log.Error(err)
		return
	}
	for _, producer := range producers {
		if producer.ProducerThingClass != "point" { // TODO: CreateProducerHistory for ProducerThingClass == "schedule"
			continue
		}
		timestamp, _ := time.Parse("2006-01-02 15:04:05+00:00", producer.Timestamp)
		if timestamp.IsZero() || currentDate.Sub(timestamp).Seconds() >= float64(*producer.HistoryInterval*60) {
			latestPH := new(model.ProducerHistory)
			// minutes are placing such a way if 15, then it will store values on 0, 15, 30, 45
			minute := (currentDate.Minute() / *producer.HistoryInterval) * *producer.HistoryInterval
			latestPH.ProducerUUID = producer.UUID
			latestPH.Timestamp = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				currentDate.Hour(), minute, 0, 0, time.UTC)
			latestPH.PresentValue = producer.PresentValue
			producerHistories = append(producerHistories, latestPH)
		}
	}
	if len(producerHistories) > 0 {
		_, err := h.DB.CreateBulkProducerHistory(producerHistories)
		if err != nil {
			log.Error(err)
		}
	}
	log.Debug("Finished create interval histories process")
}
