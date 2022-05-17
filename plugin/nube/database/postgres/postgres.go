package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/plugin/nube/database/postgres/pgmodel"
	postgresql "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func New(postgresSetting *PostgresSetting) (*PostgresSetting, error) {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		postgresSetting.Host, postgresSetting.User, postgresSetting.Password, postgresSetting.DbName,
		postgresSetting.Port, postgresSetting.SslMode)
	db, err := gorm.Open(postgresql.Open(dns), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, err
	}
	postgresConnectionInstance = &PostgresConnection{
		db: db,
	}
	postgresSetting.postgresConnectionInstance = postgresConnectionInstance
	if err := autoMigrate(db); err != nil {
		return nil, err
	}
	return postgresSetting, nil
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
	for _, s := range interfaces {
		if err := db.AutoMigrate(s); err != nil {
			return err
		}
	}
	return nil
}

func (ps PostgresSetting) WriteToPostgresDb(value interface{}) {
	ps.postgresConnectionInstance.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(value)
}
