package history

import (
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

func (h *History) InitProducerHistoryCleaner(frequency int, dataPersistingHours int) {
	h.cron = gocron.NewScheduler(time.UTC)
	_, _ = h.cron.Every(frequency).Tag("ProducerHistoryCleaner").Do(h.cleanProducerHistory, dataPersistingHours)
	h.cron.StartAsync()
}

func (h *History) cleanProducerHistory(dataPersistingHours int) {
	log.Info("Producer history cleaner has is been called...")
	persistenceTs := time.Now().UTC().Add(-time.Hour * time.Duration(dataPersistingHours)).Format(time.RFC3339Nano)
	_, err := h.DB.DeleteProducerHistoriesBeforeTimestamp(persistenceTs)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Finished producer history cleaning process")
}
