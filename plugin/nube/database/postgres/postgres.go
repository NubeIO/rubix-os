package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/database/postgres/pgmodel"
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
	}
	if (db.Migrator().HasConstraint(&pgmodel.FlowNetworkClone{}, "flow_network_clones_global_uuid_key")) {
		_ = db.Migrator().DropConstraint(&pgmodel.FlowNetworkClone{}, "flow_network_clones_global_uuid_key")
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
		return ps.postgresConnectionInstance.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(value).Error
	}
	return nil
}
