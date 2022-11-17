package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/database/postgres/pgmodel"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	postgresql "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strconv"
	"strings"
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

func (ps *PostgresSetting) New() error {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
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
		pgmodel.FlowNetworkClone{},
		pgmodel.StreamClone{},
		pgmodel.Consumer{},
		pgmodel.Writer{},
		pgmodel.Network{},
		pgmodel.Device{},
		pgmodel.Point{},
		pgmodel.Tag{},
		pgmodel.History{},
		pgmodel.NetworkMetaTag{},
		pgmodel.DeviceMetaTag{},
		pgmodel.PointMetaTag{},
	}
	if (db.Migrator().HasConstraint(&pgmodel.FlowNetworkClone{}, "flow_network_clones_global_uuid_key")) {
		_ = db.Migrator().DropConstraint(&pgmodel.FlowNetworkClone{}, "flow_network_clones_global_uuid_key")
	}
	if (db.Migrator().HasIndex(&pgmodel.Consumer{}, "idx_consumers_producer_uuid")) {
		_ = db.Migrator().DropIndex(&pgmodel.Consumer{}, "idx_consumers_producer_uuid")
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

func (ps PostgresSetting) WriteToPostgresDb(value interface{}) error {
	if reflect.ValueOf(value).Len() > 0 {
		return ps.postgresConnectionInstance.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(value, 1000).Error
	}
	return nil
}

func (ps PostgresSetting) GetHistories(args Args) ([]*pgmodel.HistoryData, error) {
	var historyDataModel []*pgmodel.HistoryData
	query, err := ps.buildHistoryQuery(args)
	if err != nil {
		return nil, err
	}
	query.Find(&historyDataModel)
	if query.Error != nil {
		return nil, errors.New("invalid filter")
	}
	return historyDataModel, nil

}

func (ps PostgresSetting) buildHistoryQuery(args Args) (*gorm.DB, error) {
	filterQuery, hasFnc, err := buildFilterQuery(args.Filter)
	if err != nil {
		return nil, err
	}
	selectQuery := buildSelectQuery(hasFnc)
	query := ps.postgresConnectionInstance.db
	query = query.Table("histories").
		Select(selectQuery).
		Joins("INNER JOIN consumers ON consumers.producer_uuid = histories.uuid").
		Joins("INNER JOIN writers ON writers.consumer_uuid = consumers.uuid").
		Joins("INNER JOIN points ON points.uuid = writers.writer_thing_uuid").
		Joins("INNER JOIN devices ON devices.uuid = points.device_uuid").
		Joins("INNER JOIN networks ON networks.uuid = devices.network_uuid")
	if hasFnc {
		query = query.Joins("INNER JOIN stream_clones ON stream_clones.uuid = consumers.stream_clone_uuid").
			Joins("INNER JOIN flow_network_clones ON flow_network_clones.uuid = stream_clones.flow_network_clone_uuid")
	}
	if args.GroupLimit != nil {
		groupLimitQuery := fmt.Sprintf("INNER JOIN (SELECT *,row_number FROM (SELECT *,ROW_NUMBER() OVER "+
			"(PARTITION BY UUID) AS row_number FROM histories ORDER BY timestamp DESC) _ WHERE row_number <= %s) AS "+
			"group_limits ON histories.id = group_limits.id AND histories.uuid = group_limits.uuid AND "+
			"histories.value = group_limits.value AND histories.timestamp = group_limits.timestamp", *args.GroupLimit)
		query = query.Joins(groupLimitQuery)
	}
	if args.Filter != nil {
		query = query.Where(filterQuery)
	}
	if args.Limit != nil {
		limit, err := strconv.Atoi(*args.Limit)
		if err == nil {
			query.Limit(limit)
		}
	}
	if args.Offset != nil {
		offset, err := strconv.Atoi(*args.Offset)
		if err == nil {
			query.Offset(offset)
		}
	}
	if args.OrderBy != nil || args.Order != nil {
		order := "DESC"
		orderBy := "timestamp"
		if args.Order != nil && strings.ToUpper(strings.TrimSpace(*args.Order)) == "ASC" {
			order = "ASC"
		}
		if args.OrderBy != nil {
			orderBy = *args.OrderBy
		}
		query.Order(fmt.Sprintf("%s %s", orderBy, order))
	}
	return query, nil
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

func (ps PostgresSetting) updateTags(model, tags interface{}) error {
	return ps.postgresConnectionInstance.db.Model(model).Association("Tags").Replace(tags)
}
