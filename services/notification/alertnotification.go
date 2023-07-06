package notification

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func (h *Notification) InitAlertNotification(frequency int) {
	h.cron = gocron.NewScheduler(time.UTC)
	_, _ = h.cron.Every(frequency).Tag("AlertNotification").Do(h.sendAlertNotification)
	h.cron.StartAsync()
}

func (h *Notification) sendAlertNotification() {
	log.Info("Send alert notification has is been called...")
	alerts, err := h.DB.GetAlerts(api.Args{Target: nstring.New("mobile"), Notified: boolean.NewFalse()})
	if err != nil {
		return
	}
	var alertsUUIDs []*string
	wg := &sync.WaitGroup{}
	for _, alert := range alerts {
		wg.Add(1)
		alertsUUIDs = append(alertsUUIDs, nstring.New(alert.UUID))
		go func(alert *model.Alert) {
			defer wg.Done()
			data := map[string]interface{}{
				"to": "",
				"notification": map[string]string{
					"title": alert.Title,
					"body":  alert.Body,
				},
				"content_available": true,
				"priority":          "high",
			}
			members, _ := h.DB.GetMembersByHostUUID(alert.HostUUID)
			for _, member := range members {
				h.DB.SendNotificationByMemberUUID(member.UUID, data)
			}
		}(alert)
	}
	wg.Wait()
	h.DB.UpdateAlertsNotified(alertsUUIDs, boolean.NewTrue())
	log.Info("Finished send alert notification process")
}
