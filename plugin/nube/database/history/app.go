package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/src/client"
	"github.com/NubeIO/rubix-os/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) syncHistory() (bool, error) {
	log.Info("History sync has is been called...")
	hosts, err := inst.db.GetHosts(api.Args{})
	if err != nil {
		return false, err
	}
	var histories []*model.History
	var historyLogs []*model.HistoryLog
	for _, host := range hosts {
		cloneEdge := true
		hisLog, err := inst.db.GetHistoryLogByHostUUID(host.UUID)
		if err != nil {
			continue
		}
		cli := client.NewClient(host.IP, host.Port, host.ExternalToken)
		pHistories, err := cli.GetPointHistoriesForSync(hisLog.LastSyncID, hisLog.Timestamp)
		if err != nil {
			continue
		}
		for k, h := range *pHistories {
			if cloneEdge {
				point, _ := inst.db.GetOnePointByArgs(args.Args{SourceUUID: nstring.New(h.PointUUID)})
				if point == nil {
					err = inst.db.CloneEdge(host)
					cloneEdge = err != nil
				}
			}
			history := model.History{
				ID:        h.ID,
				PointUUID: h.PointUUID,
				HostUUID:  host.UUID,
				Value:     h.Value,
				Timestamp: h.Timestamp,
			}
			histories = append(histories, &history)
			if k == len(*pHistories)-1 { // Update History Log
				hisLog.HostUUID = host.UUID
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
