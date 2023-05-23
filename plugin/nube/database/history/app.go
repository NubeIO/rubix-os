package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/src/client"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) syncHistory() (bool, error) {
	log.Info("History sync has is been called...")
	fnClones, err := inst.db.GetFlowNetworkClones(api.Args{})
	if err != nil {
		return false, err
	}
	var histories []*model.History
	var historyLogs []*model.HistoryLog
	for _, fnc := range fnClones {
		hisLog, err := inst.db.GetHistoryLogByFlowNetworkCloneUUID(fnc.UUID)
		if err != nil {
			continue
		}
		cli := client.NewFlowClientCliFromFNC(fnc)
		pHistories, err := cli.GetProducerHistoriesPointsForSync(hisLog.LastSyncID, hisLog.Timestamp)
		if err != nil {
			continue
		}
		for k, h := range *pHistories {
			h := h // more: https://medium.com/swlh/use-pointer-of-for-range-loop-variable-in-go-3d3481f7ffc9
			histories = append(histories, &h)
			if k == len(*pHistories)-1 { // Update History Log
				hisLog.FlowNetworkCloneUUID = fnc.UUID
				hisLog.LastSyncID = h.ID
				hisLog.Timestamp = h.Timestamp
				historyLogs = append(historyLogs, hisLog)
			}
		}
	}
	if len(histories) > 0 {
		_, err = inst.db.CreateBulkHistory(histories)
		if err != nil {
			return false, err
		}
		if len(historyLogs) > 0 {
			_, err = inst.db.UpdateBulkHistoryLogs(historyLogs)
			if err != nil {
				return false, err
			}
		}
	}
	log.Info(fmt.Sprintf("Stored %v rows on %v", len(histories), path))
	return true, nil
}
