package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/database/postgres/pgmodel"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) initializePostgresSetting() *PostgresSetting {
	postgresConnection := inst.config.Postgres
	postgresSetting := new(PostgresSetting)
	postgresSetting.Host = postgresConnection.Host
	postgresSetting.User = postgresConnection.User
	postgresSetting.Password = postgresConnection.Password
	postgresSetting.DbName = postgresConnection.DbName
	postgresSetting.Port = postgresConnection.Port
	postgresSetting.SslMode = postgresConnection.SslMode
	return postgresSetting
}

func (inst *Instance) syncPostgres(postgresSetting *PostgresSetting) (bool, error) {
	log.Info("Postgres sync has is been called...")
	if postgresSetting == nil {
		err := "postgres: Postgres sync failure: no any valid Postgres connection"
		log.Warn(err)
		return false, errors.New(err)
	}
	lastSyncId := 0
	lastSyncHistoryPostgresLog, err := inst.db.GetLastSyncHistoryPostgresLog()
	if lastSyncHistoryPostgresLog != nil {
		lastSyncId = lastSyncHistoryPostgresLog.ID
	}
	if err != nil {
		log.Warn(err)
		return false, err
	}
	histories, err := inst.db.GetHistoriesForPostgresSync(lastSyncId)
	if err != nil {
		log.Warn(err)
		return false, err
	}
	if len(histories) > 0 {
		// bulk write to postgres
		flowNetworkClones, _ := inst.db.GetFlowNetworkClones(api.Args{WithStreamClones: true, WithConsumers: true,
			WithWriters: true, WithTags: true})
		mFlowNetworkClones, err := json.Marshal(flowNetworkClones)
		if err != nil {
			log.Error(err)
			return false, err
		}
		var flowNetworkClonesModel []*pgmodel.FlowNetworkClone
		if err = json.Unmarshal(mFlowNetworkClones, &flowNetworkClonesModel); err != nil {
			log.Error(err)
			return false, err
		}
		postgresSetting.WriteToPostgresDb(flowNetworkClonesModel)

		networks, _ := inst.db.GetNetworks(api.Args{WithTags: true, WithDevices: true, WithPoints: true})
		mNetworks, err := json.Marshal(networks)
		if err != nil {
			log.Error(err)
		}
		var networksModel []*pgmodel.Network
		if err = json.Unmarshal(mNetworks, &networksModel); err != nil {
			log.Error(err)
			return false, err
		}
		postgresSetting.WriteToPostgresDb(networksModel)

		postgresSetting.WriteToPostgresDb(histories)
		lastHistory := histories[len(histories)-1]
		historyPostgresLog := &model.HistoryPostgresLog{
			ID:        lastHistory.ID,
			UUID:      lastHistory.UUID,
			Value:     lastHistory.Value,
			Timestamp: lastHistory.Timestamp,
		}
		_, _ = inst.db.UpdateHistoryPostgresLog(historyPostgresLog)
		log.Info(fmt.Sprintf("postgres: Stored %v new records", len(histories)))
	} else {
		log.Info("postgres: Nothing to store, no new records")
	}
	return true, nil
}
