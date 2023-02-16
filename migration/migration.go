package migration

import (
	"github.com/NubeIO/flow-framework/migration/versions"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if (db.Migrator().HasColumn(&model.History{}, "idx_histories_id") ||
		!db.Migrator().HasColumn(&model.History{}, "history_id")) {
		// We drop old table and create new one with composite primary_key
		// Update of primary_key doesn't support in GORM
		// Fore more info: https://github.com/go-gorm/gorm/issues/4742
		err := db.Migrator().DropTable(&model.History{})
		log.Error(err)
	}
	if db.Migrator().HasIndex(&model.Device{}, "name_network_composite_index") {
		err := db.Migrator().DropIndex(&model.Device{}, "name_network_composite_index")
		log.Error(err)
	}
	if db.Migrator().HasIndex(&model.FlowNetwork{}, "ip_port_composite_index") {
		err := db.Migrator().DropIndex(&model.FlowNetwork{}, "ip_port_composite_index")
		log.Error(err)
	}
	if db.Migrator().HasIndex(&model.FlowNetworkClone{}, "ip_port_clone_composite_index") {
		err := db.Migrator().DropIndex(&model.FlowNetworkClone{}, "ip_port_clone_composite_index")
		log.Error(err)
	}
	if db.Migrator().HasIndex(&model.History{}, "id_uuid_value_timestamp_composite_index") {
		err := db.Migrator().DropIndex(&model.History{}, "id_uuid_value_timestamp_composite_index")
		log.Error(err)
	}
	if db.Migrator().HasIndex(&model.Point{}, "name_device_composite_index") {
		err := db.Migrator().DropIndex(&model.Point{}, "name_device_composite_index")
		log.Error(err)
	}
	if db.Migrator().HasColumn(&model.Point{}, "history_interval") {
		columnTypes, _ := db.Migrator().ColumnTypes(&model.Point{})
		for _, columnType := range columnTypes {
			if columnType.Name() == "history_interval" && columnType.DatabaseTypeName() == "real" {
				err := db.Migrator().DropColumn(&model.Point{}, "history_interval")
				log.Error(err)
				break
			}
		}
	}
	interfaces := versions.GetInitInterfaces()
	for _, s := range interfaces {
		if err := db.AutoMigrate(s); err != nil {
			return err
		}
	}
	return nil
}
