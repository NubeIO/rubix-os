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
	interfaces := versions.GetInitInterfaces()
	for _, s := range interfaces {
		if err := db.AutoMigrate(s); err != nil {
			return err
		}
	}
	return nil
}
