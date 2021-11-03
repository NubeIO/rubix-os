package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) syncHistory() (bool, error) {
	log.Info("History sync has is been called")
	fnClones, err := i.db.GetFlowNetworkClones()
	if err != nil {
		return false, err
	}
	var histories []*model.History
	for _, fnClone := range fnClones {
		sClones := fnClone.StreamClones
		for _, sClone := range sClones {
			url := fmt.Sprintf("%s:%v/%s", *fnClone.FlowIP, *fnClone.FlowPort, sClone.UUID)
			sHistories, _ := getProducerHistories(url)
			histories = append(histories, sHistories...)
			//example for testing only
			//histories = append(histories, &model.History{
			//	UUID:      sClone.UUID,
			//	ID:        0,
			//	Timestamp: time.Now(),
			//	Value:     rand.Float64(),
			//})
		}
	}
	if len(histories) > 0 {
		//save histories
		_, err = i.db.CreateBulkHistory(histories)
		if err != nil {
			return false, err
		}
	}
	log.Info(fmt.Sprintf("Stored %v rows on %v", len(histories), path))
	return true, nil
}
