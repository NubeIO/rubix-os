package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) syncInflux() (bool, error) {
	log.Info("InfluxDB sync has is been called...")
	influxConnections := i.config.Influx
	if len(influxConnections) == 0 {
		log.Warn("InfluxDB sync failure: integration details missing")
		return true, nil
	}
	histories, err := i.db.GetHistoriesForSync()
	if err != nil {
		return false, err
	}
	lastSyncId := 0
	producerUuid := ""
	var historyTags []*model.HistoryInfluxTag

	for _, history := range histories {
		if lastSyncId < history.ID {
			lastSyncId = history.ID
		}
		if producerUuid != history.UUID {
			producerUuid = history.UUID
			historyTags, err = i.db.GetHistoryInfluxTags(producerUuid)
			if err != nil {
				log.Error(fmt.Sprintf("We unable to get the producer_uuid = %s details!", producerUuid))
			}
		}
		for _, historyTag := range historyTags {
			tags := tagsHistory(historyTag)
			fields := fieldsHistory(history)
			for _, influx := range influxConnections {
				influxSetting := new(InfluxSetting)
				schema := "http"
				if influx.Port == 443 {
					schema = "https"
				}
				influxSetting.ServerURL = fmt.Sprintf("%s://%s:%d", schema, influx.Host, influx.Port)
				if influx.Token == nil {
					log.Warn("Token is null, please update it")
					continue
				}
				influxSetting.AuthToken = *influx.Token
				influxSetting.Org = influx.Org
				influxSetting.Bucket = influx.Bucket
				influxSetting.Measurement = influx.Measurement
				isc := New(influxSetting)
				isc.WriteHistories(tags, fields, history.Timestamp)
			}
		}
	}
	historyCount := len(histories)
	if historyCount > 0 {
		_, err := i.db.UpdateHistoryInfluxLogLastSyncId(lastSyncId)
		if err != nil {
			return false, err
		}
		log.Info(fmt.Sprintf("Stored %v rows on %v", historyCount, path))
	} else {
		log.Info("Nothing to store, no new records")
	}
	return true, nil
}
