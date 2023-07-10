package main

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/plugin/nube/database/postgres/pgmodel"
	log "github.com/sirupsen/logrus"
	postgresql "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"time"
)

var postgresConnectionInstance *PostgresConnection

type PostgresSetting struct {
	Host                       string
	User                       string
	Password                   string
	DbName                     string
	Port                       int
	SslMode                    string
	postgresConnectionInstance *PostgresConnection
}

type PostgresConnection struct {
	db *gorm.DB
}

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

func (inst *Instance) initializePostgresDBConnection() (bool, error) {
	inst.cpsDebugMsg("initializePostgresDBConnection")
	if postgresSetting.postgresConnectionInstance.db == nil {
		err := postgresSetting.New()
		if err != nil {
			inst.cpsErrorMsg(err)
			return false, err
		}
	}
	return true, nil
}

func (ps *PostgresSetting) New() error {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		ps.Host, ps.User, ps.Password, ps.DbName, ps.Port, ps.SslMode)
	ps.postgresConnectionInstance = &PostgresConnection{
		db: nil,
	}
	db, err := connectDb(dns)
	if err != nil {
		return err
	}
	ps.postgresConnectionInstance = &PostgresConnection{
		db: db,
	}
	if err := autoMigrate(db); err != nil {
		return err
	}
	return nil
}

func autoMigrate(db *gorm.DB) error {
	interfaces := []interface{}{
		Site{},
		Threshold{},
	}

	for _, s := range interfaces {
		if err := db.AutoMigrate(s); err != nil {
			return err
		}
	}
	return nil
}

func connectDb(dns string) (*gorm.DB, error) {
	return gorm.Open(postgresql.Open(dns), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
}

func (inst *Instance) syncNetDevPntsAndTags() (bool, error) {
	/*
		points, err := inst.db.GetPointsTable()
		if err != nil {
			return false, err
		}
		for i, point := range points {
			if point != nil {
				inst.cpsDebugMsg(fmt.Sprintf("GetPointsTable() point %v: %+v", i, *point))
			}
		}
		return true, nil

	*/

	inst.cpsDebugMsg("syncNetDevPntsAndTags() has been called...")
	if postgresSetting.postgresConnectionInstance.db == nil {
		err := postgresSetting.New()
		if err != nil {
			log.Warn(err)
			return false, err
		}
	}

	if err := inst.createPointsBulk(); err != nil {
		inst.cpsErrorMsg("syncNetDevPntsAndTags() error: ", err)
		return false, err
	}
	return true, nil
}

func (inst *Instance) createPointsBulk() error {
	points, err := inst.db.GetPointsForPostgresSync()
	if err != nil {
		return err
	}
	for i, point := range points {
		if point != nil {
			inst.cpsDebugMsg(fmt.Sprintf("syncNetDevPntsAndTags() point %v: %+v", i, *point))
		}
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

func (ps PostgresSetting) WriteToPostgresDb(value interface{}) error {
	if reflect.ValueOf(value).Len() > 0 {
		return ps.postgresConnectionInstance.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(value, 1000).Error
	}
	return nil
}

func convertData(data interface{}, v interface{}) error {
	mData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(mData, &v); err != nil {
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

func (ps PostgresSetting) DeleteDeletedNetworkMetaTags(metaTags []*model.NetworkMetaTag) error {
	if len(metaTags) == 0 {
		return ps.postgresConnectionInstance.db.Where("true").Delete(pgmodel.NetworkMetaTag{}).Error
	}
	notIn := make([][]interface{}, len(metaTags))
	for i, metaTag := range metaTags {
		ni := make([]interface{}, 2)
		ni[0] = metaTag.NetworkUUID
		ni[1] = metaTag.Key
		notIn[i] = ni
	}
	return ps.postgresConnectionInstance.db.Where("(network_uuid,key) NOT IN ?", notIn).
		Delete(pgmodel.NetworkMetaTag{}).Error
}

func (ps PostgresSetting) DeleteDeletedDeviceMetaTags(metaTags []*model.DeviceMetaTag) error {
	if len(metaTags) == 0 {
		return ps.postgresConnectionInstance.db.Where("true").Delete(pgmodel.DeviceMetaTag{}).Error
	}
	notIn := make([][]interface{}, len(metaTags))
	for i, metaTag := range metaTags {
		ni := make([]interface{}, 2)
		ni[0] = metaTag.DeviceUUID
		ni[1] = metaTag.Key
		notIn[i] = ni
	}
	return ps.postgresConnectionInstance.db.Where("(device_uuid,key) NOT IN ?", notIn).
		Delete(pgmodel.DeviceMetaTag{}).Error
}

func (ps PostgresSetting) DeleteDeletedPointMetaTags(metaTags []*model.PointMetaTag) error {
	if len(metaTags) == 0 {
		return ps.postgresConnectionInstance.db.Where("true").Delete(pgmodel.PointMetaTag{}).Error
	}
	notIn := make([][]interface{}, len(metaTags))
	for i, metaTag := range metaTags {
		ni := make([]interface{}, 2)
		ni[0] = metaTag.PointUUID
		ni[1] = metaTag.Key
		notIn[i] = ni
	}
	return ps.postgresConnectionInstance.db.Where("(point_uuid,key) NOT IN ?", notIn).
		Delete(pgmodel.PointMetaTag{}).Error
}

func (ps PostgresSetting) DeleteDeletedNetworkTags(tags []*pgmodel.NetworkTag) error {
	if len(tags) == 0 {
		return ps.postgresConnectionInstance.db.Where("true").Delete(pgmodel.NetworkTag{}).Error
	}
	notIn := make([][]interface{}, len(tags))
	for i, tag := range tags {
		ni := make([]interface{}, 2)
		ni[0] = tag.NetworkUUID
		ni[1] = tag.Tag
		notIn[i] = ni
	}
	return ps.postgresConnectionInstance.db.Where("(network_uuid,tag) NOT IN ?", notIn).
		Delete(pgmodel.NetworkTag{}).Error
}

func (ps PostgresSetting) DeleteDeletedDeviceTags(tags []*pgmodel.DeviceTag) error {
	if len(tags) == 0 {
		return ps.postgresConnectionInstance.db.Where("true").Delete(pgmodel.DeviceTag{}).Error
	}
	notIn := make([][]interface{}, len(tags))
	for i, tag := range tags {
		ni := make([]interface{}, 2)
		ni[0] = tag.DeviceUUID
		ni[1] = tag.Tag
		notIn[i] = ni
	}
	return ps.postgresConnectionInstance.db.Where("(device_uuid,tag) NOT IN ?", notIn).
		Delete(pgmodel.DeviceTag{}).Error
}

func (ps PostgresSetting) DeleteDeletedPointTags(tags []*pgmodel.PointTag) error {
	if len(tags) == 0 {
		return ps.postgresConnectionInstance.db.Where("true").Delete(pgmodel.PointTag{}).Error
	}
	notIn := make([][]interface{}, len(tags))
	for i, tag := range tags {
		ni := make([]interface{}, 2)
		ni[0] = tag.PointUUID
		ni[1] = tag.Tag
		notIn[i] = ni
	}
	return ps.postgresConnectionInstance.db.Where("(point_uuid,tag) NOT IN ?", notIn).
		Delete(pgmodel.PointTag{}).Error
}

func (inst *Instance) SendHistoriesToPostgres(histories []*pgmodel.History) (bool, error) {
	inst.cpsDebugMsg("SendHistoriesToPostgres()")
	if postgresSetting.postgresConnectionInstance.db == nil {
		err := postgresSetting.New()
		if err != nil {
			log.Warn(err)
			return false, err
		}
	}

	// TODO: deal with the last sync of processed data
	if len(histories) > 0 {
		if err := postgresSetting.WriteToPostgresDb(histories); err != nil {
			inst.cpsErrorMsg("SendHistoriesToPostgres() error:", err)
			return false, err
		}
		inst.cpsDebugMsg(fmt.Sprintf("SendHistoriesToPostgres(): Stored %v new records", len(histories)))
	} else {
		log.Info("postgres: Nothing to store, no new records")
	}
	return true, nil
}
