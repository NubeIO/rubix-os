package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/plugin/nube/database/postgres/pgmodel"
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
	// _, err := inst.db.GetHistoryPostgresLogLastSyncHistoryId()
	if err != nil {
		log.Warn(err)
		return false, err
	}

	histories, err := inst.db.GetHistoriesForPostgresSync(lastSyncId)
	if err != nil {
		log.Warn(err)
		return false, err
	}

	/*
		// TODO: DELETE ME.  this section creates example histories from csv files.
		var histories []*model.History
		// Open the CSV file
		filePath := "/home/marc/Documents/Nube/CPS/Development/Data_Processing/Example_Sensor_Data/Door_Sensor_Data_Example_Upload_3.csv"
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Failed to open CSV file: %v", err)
		}
		defer file.Close()
		// Read and process the CSV data
		reader := csv.NewReader(file)
		reader.FieldsPerRecord = 5
		first := true
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if first {
				first = false
				continue
			}
			if err != nil {
				log.Printf("Failed to read CSV record: %v", err)
				continue
			}
			log.Printf("record: %+v", record)

			// Parse CSV record and create a History instance
			history := model.History{}
			history.ID, _ = strconv.Atoi(record[0])
			history.PointUUID = record[1]
			history.HostUUID = record[2]
			value, _ := strconv.ParseFloat(record[3], 64)
			history.Value = &value
			timestamp, _ := time.Parse(time.RFC3339, record[4])
			history.Timestamp = timestamp

			histories = append(histories, &history)
		}

	*/

	if len(histories) > 0 {
		if !inst.config.Job.DisableTagSync {
			if err = inst.createPointsBulk(); err != nil {
				log.Error(err)
				return false, err
			}
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
			PointUUID: lastHistory.PointUUID,
			HostUUID:  lastHistory.HostUUID,
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

func (inst *Instance) getHistories(args Args) ([]*pgmodel.HistoryDataResponse, error) {
	if postgresSetting.postgresConnectionInstance.db == nil {
		err := postgresSetting.New()
		if err != nil {
			log.Warn(err)
			return nil, err
		}
	}
	histories, err := postgresSetting.GetHistories(args)
	if err != nil {
		return nil, err
	}
	historiesResponse := convertHistoryDataToResponse(histories)
	return historiesResponse, nil
}

func (inst *Instance) createPointsBulk() error {
	points, err := inst.db.GetPointsForPostgresSync()
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
	if err = inst.createTags(); err != nil {
		return err
	}
	if err = inst.createMetaTags(); err != nil {
		return err
	}
	return nil
}

func (inst *Instance) createTags() error {
	networkTags, err := inst.db.GetNetworksTagsForPostgresSync()
	if err != nil {
		return err
	}
	var networkTagsModel []*pgmodel.NetworkTag
	if err = convertData(networkTags, &networkTagsModel); err != nil {
		return err
	}
	if err := postgresSetting.DeleteDeletedNetworkTags(networkTagsModel); err != nil {
		return err
	}
	if len(networkTagsModel) > 0 {
		if err := postgresSetting.WriteToPostgresDb(networkTagsModel); err != nil {
			return err
		}
	}

	deviceTags, err := inst.db.GetDevicesTagsForPostgresSync()
	if err != nil {
		return err
	}
	var deviceTagsModel []*pgmodel.DeviceTag
	if err = convertData(deviceTags, &deviceTagsModel); err != nil {
		return err
	}
	if err := postgresSetting.DeleteDeletedDeviceTags(deviceTagsModel); err != nil {
		return err
	}
	if len(deviceTagsModel) > 0 {
		if err := postgresSetting.WriteToPostgresDb(deviceTagsModel); err != nil {
			return err
		}
	}

	pointTags, err := inst.db.GetPointsTagsForPostgresSync()
	if err != nil {
		return err
	}
	var pointTagsModel []*pgmodel.PointTag
	if err = convertData(pointTags, &pointTagsModel); err != nil {
		return err
	}
	if err := postgresSetting.DeleteDeletedPointTags(pointTagsModel); err != nil {
		return err
	}
	if len(pointTagsModel) > 0 {
		if err := postgresSetting.WriteToPostgresDb(pointTagsModel); err != nil {
			return err
		}
	}
	return nil
}

func (inst *Instance) createMetaTags() error {
	networkMetaTags, err := inst.db.GetNetworksMetaTagsForPostgresSync()
	if err != nil {
		return err
	}
	if err := postgresSetting.DeleteDeletedNetworkMetaTags(networkMetaTags); err != nil {
		return err
	}
	if len(networkMetaTags) > 0 {
		if err := postgresSetting.WriteToPostgresDb(networkMetaTags); err != nil {
			return err
		}
	}

	deviceMetaTags, err := inst.db.GetDevicesMetaTagsForPostgresSync()
	if err != nil {
		return err
	}
	if err := postgresSetting.DeleteDeletedDeviceMetaTags(deviceMetaTags); err != nil {
		return err
	}
	if len(deviceMetaTags) > 0 {
		if err := postgresSetting.WriteToPostgresDb(deviceMetaTags); err != nil {
			return err
		}
	}

	pointMetaTags, err := inst.db.GetPointsMetaTagsForPostgresSync()
	if err != nil {
		return err
	}
	if err := postgresSetting.DeleteDeletedPointMetaTags(pointMetaTags); err != nil {
		return err
	}
	if len(pointMetaTags) > 0 {
		if err := postgresSetting.WriteToPostgresDb(pointMetaTags); err != nil {
			return err
		}
	}
	return nil
}
