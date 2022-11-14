package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/database/postgres/pgmodel"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

var postgresSetting *PostgresSetting

func (inst *Instance) initializePostgresSetting() {
	postgresConnection := inst.config.Postgres
	if postgresSetting == nil {
		postgresSetting = new(PostgresSetting)
	}
	postgresSetting.Host = postgresConnection.Host
	postgresSetting.User = postgresConnection.User
	postgresSetting.Password = postgresConnection.Password
	postgresSetting.DbName = postgresConnection.DbName
	postgresSetting.Port = postgresConnection.Port
	postgresSetting.SslMode = postgresConnection.SslMode
	postgresSetting.postgresConnectionInstance = &PostgresConnection{
		db: nil,
	}
}

func (inst *Instance) syncPostgres() (bool, error) {
	log.Info("postgres sync has been called...")
	if postgresSetting.postgresConnectionInstance.db == nil {
		err := postgresSetting.New()
		if err != nil {
			log.Warn(err)
			return false, err
		}
	}
	lastSyncId, err := inst.db.GetHistoryPostgresLogLastSyncHistoryId()
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
		if err = inst.createFlowNetworkCloneBulk(); err != nil {
			log.Error(err)
			return false, err
		}
		if err = inst.createStreamCloneBulk(); err != nil {
			log.Error(err)
			return false, err
		}
		if err = inst.createConsumerBulk(); err != nil {
			log.Error(err)
			return false, err
		}
		if err = inst.createWriterBulk(); err != nil {
			log.Error(err)
			return false, err
		}
		if err = inst.createNetworkBulk(); err != nil {
			log.Error(err)
			return false, err
		}
		if err = inst.createDeviceBulk(); err != nil {
			log.Error(err)
			return false, err
		}
		if err = inst.createPointBulk(); err != nil {
			log.Error(err)
			return false, err
		}

		var historiesModel []*pgmodel.History
		if err = convertData(histories, &historiesModel); err != nil {
			return false, err
		}
		if err = postgresSetting.WriteToPostgresDb(historiesModel); err != nil {
			log.Error(err)
			return false, err
		}
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

func (inst *Instance) getHistories(args Args) ([]*pgmodel.HistoryData, error) {
	if postgresSetting.postgresConnectionInstance.db == nil {
		err := postgresSetting.New()
		if err != nil {
			log.Warn(err)
			return nil, err
		}
	}
	return postgresSetting.GetHistories(args)
}

func (inst *Instance) createFlowNetworkCloneBulk() error {
	flowNetworkClones, err := inst.db.GetFlowNetworkClones(api.Args{WithTags: true})
	if err != nil {
		return err
	}
	var flowNetworkClonesModel []*pgmodel.FlowNetworkClone
	if err = convertData(flowNetworkClones, &flowNetworkClonesModel); err != nil {
		return err
	}
	return postgresSetting.WriteToPostgresDb(flowNetworkClonesModel)
}

func (inst *Instance) createStreamCloneBulk() error {
	streamClones, err := inst.db.GetStreamClones(api.Args{WithTags: true})
	if err != nil {
		return err
	}
	var streamClonesModel []*pgmodel.StreamClone
	if err = convertData(streamClones, &streamClonesModel); err != nil {
		return err
	}
	if err = postgresSetting.WriteToPostgresDb(streamClonesModel); err != nil {
		return err
	}
	// tags
	for _, streamCloneModel := range streamClonesModel {
		if err = postgresSetting.updateTags(streamCloneModel, streamCloneModel.Tags); err != nil {
			return err
		}
	}
	return nil
}

func (inst *Instance) createConsumerBulk() error {
	consumers, err := inst.db.GetConsumers(api.Args{WithTags: true})
	if err != nil {
		return err
	}
	var consumersModel []*pgmodel.Consumer
	if err = convertData(consumers, &consumersModel); err != nil {
		return err
	}
	if err = postgresSetting.WriteToPostgresDb(consumersModel); err != nil {
		return err
	}
	// tags
	for _, consumerModel := range consumersModel {
		if err = postgresSetting.updateTags(consumerModel, consumerModel.Tags); err != nil {
			return err
		}
	}
	return nil
}

func (inst *Instance) createWriterBulk() error {
	writers, err := inst.db.GetWriters(api.Args{})
	if err != nil {
		return err
	}
	var writersModel []*pgmodel.Writer
	if err = convertData(writers, &writersModel); err != nil {
		return err
	}
	return postgresSetting.WriteToPostgresDb(writersModel)
}

func (inst *Instance) createNetworkBulk() error {
	networks, err := inst.db.GetNetworks(api.Args{WithTags: true})
	if err != nil {
		return err
	}
	var networksModel []*pgmodel.Network
	if err = convertData(networks, &networksModel); err != nil {
		return err
	}
	if err = postgresSetting.WriteToPostgresDb(networksModel); err != nil {
		return err
	}
	// tags
	for _, networkModel := range networksModel {
		if err = postgresSetting.updateTags(networkModel, networkModel.Tags); err != nil {
			return err
		}
	}
	// meta tags
	networkMetaTags, err := inst.db.GetNetworkMetaTags()
	if err != nil {
		return err
	}
	if err = postgresSetting.DeleteDeletedNetworkMetaTags(networkMetaTags); err != nil {
		return err
	}
	var networkMetaTagsModel []*pgmodel.NetworkMetaTag
	if err = convertData(networkMetaTags, &networkMetaTagsModel); err != nil {
		return err
	}
	return postgresSetting.WriteToPostgresDb(networkMetaTagsModel)
}

func (inst *Instance) createDeviceBulk() error {
	devices, err := inst.db.GetDevices(api.Args{WithTags: true})
	if err != nil {
		return err
	}
	var devicesModel []*pgmodel.Device
	if err = convertData(devices, &devicesModel); err != nil {
		return err
	}
	if err = postgresSetting.WriteToPostgresDb(devicesModel); err != nil {
		return err
	}
	// tags
	for _, deviceModel := range devicesModel {
		if err = postgresSetting.updateTags(deviceModel, deviceModel.Tags); err != nil {
			return err
		}
	}
	// meta tags
	deviceMetaTags, err := inst.db.GetDeviceMetaTags()
	if err != nil {
		return err
	}
	if err := postgresSetting.DeleteDeletedDeviceMetaTags(deviceMetaTags); err != nil {
		return err
	}
	var deviceMetaTagsModel []*pgmodel.DeviceMetaTag
	if err = convertData(deviceMetaTags, &deviceMetaTagsModel); err != nil {
		return err
	}
	return postgresSetting.WriteToPostgresDb(deviceMetaTagsModel)
}

func (inst *Instance) createPointBulk() error {
	points, err := inst.db.GetPoints(api.Args{WithTags: true})
	if err != nil {
		return err
	}
	var pointsModel []*pgmodel.Point
	if err = convertData(points, &pointsModel); err != nil {
		return err
	}
	if err = postgresSetting.WriteToPostgresDb(pointsModel); err != nil {
		return err
	}
	// tags
	for _, pointModel := range pointsModel {
		if err = postgresSetting.updateTags(pointModel, pointModel.Tags); err != nil {
			return err
		}
	}
	// meta tags
	pointMetaTags, err := inst.db.GetPointMetaTags()
	if err != nil {
		return err
	}
	if err = postgresSetting.DeleteDeletedPointMetaTags(pointMetaTags); err != nil {
		return err
	}
	var pointMetaTagsModel []*pgmodel.PointMetaTag
	if err = convertData(pointMetaTags, &pointMetaTagsModel); err != nil {
		return err
	}
	return postgresSetting.WriteToPostgresDb(pointMetaTagsModel)
}
