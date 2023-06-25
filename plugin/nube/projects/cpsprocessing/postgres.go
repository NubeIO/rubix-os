package main

import (
	"fmt"
	postgresql "gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		Site{},
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

/*

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

*/
