package main

import (
	"fmt"
	influxmodel "github.com/NubeDev/flow-framework/plugin/nube/database/influx/model"
	log "github.com/sirupsen/logrus"
)

//syncInflux
func (i *Instance) syncInflux() (bool, error) {
	log.Info("InfluxDB sync has is been called")
	integrations, err := i.db.GetEnabledIntegrationByPluginConfId(i.pluginUUID)
	if err != nil {
		return false, err
	}
	for _, integration := range integrations {
		influxSetting := new(InfluxSetting)
		influxSetting.ServerURL = integration.IP + ":" + integration.PORT
		influxSetting.AuthToken = integration.Token
		isc := New(influxSetting)
		histories, err := i.db.GetHistoriesForSync()
		if err != nil {
			return false, err
		}
		for _, h := range histories {
			var hist influxmodel.HistPayload
			hist.ID = h.ID
			hist.UUID = h.UUID
			hist.Timestamp = h.Timestamp
			hist.Value = h.Value
			isc.WriteHist(hist)
		}
		historyCount := len(histories)
		if historyCount > 0 {
			// TODO: Enju
			//_, err := i.db.UpdateHistoryLogLastSyncId(histories[historyCount-1].ID)
			if err != nil {
				return false, err
			}
			log.Info(fmt.Sprintf("Stored %v rows on %v", historyCount, path))
		} else {
			log.Info("Nothing to store, no new records")
		}
	}
	if len(integrations) == 0 {
		log.Info("InfluxDB can't be registered, integration details missing.")
	}
	return true, nil
}
