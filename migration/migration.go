package migration

import (
	"github.com/NubeIO/flow-framework/migration/versions"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	interfaces := versions.GetInitInterfaces()
	for _, s := range interfaces {
		if err := db.AutoMigrate(s); err != nil {
			return err
		}
	}
	return nil
}
