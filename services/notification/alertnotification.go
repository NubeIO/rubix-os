package notification

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func (h *Notification) InitAlertNotification(frequency, resendDuration time.Duration) {
	h.cron = gocron.NewScheduler(time.UTC)
	h.cron.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
	_, _ = h.cron.Every(frequency).Tag("AlertNotification").Do(h.sendAlertNotification, resendDuration)
	h.cron.StartAsync()
}

func (h *Notification) sendAlertNotification(resendDuration time.Duration) {
	log.Info("Send alert notification has is been called...")
	notifiedAtLt := time.Now().UTC().Add(-resendDuration).Format(time.RFC3339Nano)
	alerts, err := h.DB.GetAlertsForNotification(notifiedAtLt)
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
			uniqueDevices := map[string]string{}
			for _, member := range members {
				memberDevices, _ := h.DB.GetMemberDevicesByMemberUUID(member.UUID)
				for _, memberDevice := range memberDevices {
					uniqueDevices[memberDevice.DeviceID] = *memberDevice.DeviceName
				}
			}
			h.DB.SendNotificationByMemberUUID(uniqueDevices, data)
		}(alert)
	}
	wg.Wait()
	h.DB.UpdateAlertsNotified(alertsUUIDs, boolean.NewTrue())
	log.Info("Finished send alert notification process")
}
