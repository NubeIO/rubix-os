package history

import (
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

func (h *History) InitPointHistoryCleaner(frequency int, dataPersistingHours int) {
	h.cron = gocron.NewScheduler(time.UTC)
	_, _ = h.cron.Every(frequency).Tag("PointHistoryCleaner").Do(h.cleanPointHistory, dataPersistingHours)
	h.cron.StartAsync()
}

func (h *History) cleanPointHistory(dataPersistingHours int) {
	log.Info("Point history cleaner has is been called...")
	persistenceTs := time.Now().UTC().Add(-time.Hour * time.Duration(dataPersistingHours)).Format(time.RFC3339Nano)
	_, err := h.DB.DeletePointHistoriesBeforeTimestamp(persistenceTs)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Finished point history cleaning process")
}
