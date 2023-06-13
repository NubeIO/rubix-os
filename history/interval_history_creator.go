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
	var pointHistories []*model.PointHistory
	currentDate := time.Now().UTC()

	points, err := h.DB.GetPointsForCreateInterval()
	if err != nil {
		log.Error(err)
		return
	}
	for _, point := range points {
		timestamp, _ := time.Parse("2006-01-02 15:04:05+00:00", point.Timestamp)
		if timestamp.IsZero() || currentDate.Sub(timestamp).Seconds() >= float64(*point.HistoryInterval*60) {
			latestPH := new(model.PointHistory)
			// minutes are placing such a way if 15, then it will store values on 0, 15, 30, 45
			minute := (currentDate.Minute() / *point.HistoryInterval) * *point.HistoryInterval
			latestPH.PointUUID = point.UUID
			latestPH.Timestamp = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				currentDate.Hour(), minute, 0, 0, time.UTC)
			latestPH.Value = point.PresentValue
			pointHistories = append(pointHistories, latestPH)
		}
	}
	if len(pointHistories) > 0 {
		_, err := h.DB.CreateBulkPointHistory(pointHistories)
		if err != nil {
			log.Error(err)
		}
	}
	log.Debug("Finished create interval histories process")
}
