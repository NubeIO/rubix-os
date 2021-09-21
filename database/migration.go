package database

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "201608301400",
			Migrate: func(tx *gorm.DB) error {
				type Alert struct {
					Type     string
					Duration int
				}
				type Message struct {
					ID            uint `gorm:"AUTO_INCREMENT;primary_key;index"`
					ApplicationID uint
					Message       string `gorm:"type:text"`
					Title         string `gorm:"type:text"`
					Priority      int
					Extras        []byte
					Date          time.Time
				}
				return tx.AutoMigrate(
					&Alert{},
					&Message{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"alerts",
					"messages",
				)
			},
		},
		{
			ID: "201608301500",
			Migrate: func(tx *gorm.DB) error {
				nodes := `CREATE TABLE "nodes" (
					"uuid"	varchar(255) UNIQUE,
					"name"	text,
					"node_type"	text,
					"help"	text,
					"in1"	text,
					"in2"	text,
					"out1_value"	text,
					"out2_value"	text,
					"node_settings"	JSON,
					PRIMARY KEY("uuid")
				);`
				return tx.Exec(nodes).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"nodes",
				)
			},
		},
	})
	return m.Migrate()
}
