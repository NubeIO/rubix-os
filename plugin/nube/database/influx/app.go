package main

import (
	"encoding/json"
	"github.com/NubeDev/flow-framework/model"
	influxmodel "github.com/NubeDev/flow-framework/plugin/nube/database/influx/model"
)

//writeData to influx
func (i *Instance) writeData() (bool, error) {
	histories, err := i.db.GetProducerHistories()
	var hist influxmodel.HistPayload
	var histPri *model.Priority
	if err != nil {
		return true, nil
	}
	for _, e := range histories {
		err := json.Unmarshal(e.DataStore, &histPri)
		if err != nil {
			return false, err
		}
		hist.WriterUUID = e.WriterUUID
		hist.ThingWriterUUID = e.ThingWriterUUID
		hist.ProducerUUID = e.ProducerUUID
		hist.Timestamp = e.Timestamp
		hist.Out16 = histPri.P16
		WriteHist(hist)

	}
	return true, nil
}
