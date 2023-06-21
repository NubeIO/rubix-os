package migration

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/migration/versions"
	"github.com/NubeIO/rubix-os/schema/loraschema"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if (db.Migrator().HasColumn(&model.History{}, "idx_histories_id") ||
		!db.Migrator().HasColumn(&model.History{}, "history_id") ||
		!db.Migrator().HasColumn(&model.History{}, "point_uuid")) {
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
	if db.Migrator().HasIndex(&model.Group{}, "idx_networks_name_location_uuid") {
		err := db.Migrator().DropIndex(&model.Group{}, "idx_networks_name_location_uuid")
		log.Error(err)
	}
	if (!db.Migrator().HasColumn(&model.HistoryPostgresLog{}, "point_uuid")) {
		err := db.Migrator().DropTable(&model.HistoryPostgresLog{})
		log.Error(err)
	}

	// TODO: if we uncomment this, it will remove the Priority sub-table from point table as well
	// if db.Migrator().HasColumn(&model.Point{}, "history_interval") {
	//	columnTypes, _ := db.Migrator().ColumnTypes(&model.Point{})
	//	for _, columnType := range columnTypes {
	//		if columnType.Name() == "history_interval" && columnType.DatabaseTypeName() == "real" {
	//			err := db.Migrator().DropColumn(&model.Point{}, "history_interval")
	//			log.Error(err)
	//			break
	//		}
	//	}
	// }
	deviceModel := model.Device{CommonDevice: model.CommonDevice{Model: loraschema.DeviceModelMicroEdgeV1}}
	db.Model(&model.Device{}).
		Select("Model").
		Where("model = ? OR model = ?", "MicroEdge", "MICROEDGE").
		Updates(&deviceModel)

	interfaces := versions.GetInitInterfaces()
	for _, s := range interfaces {
		if err := db.AutoMigrate(s); err != nil {
			return err
		}
	}
	return nil
}
