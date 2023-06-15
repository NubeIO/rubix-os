package database

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/user"
	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/migration"
	"github.com/NubeIO/rubix-os/module/shared"
	"github.com/NubeIO/rubix-os/src/cachestore"
	"github.com/NubeIO/rubix-registry-go/rubixregistry"
	"os"
	"path/filepath"
	"time"

	"github.com/NubeIO/rubix-os/logger"
	"github.com/NubeIO/rubix-os/plugin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	username = "admin"
	password = "N00BWires"
)

var mkdirAll = os.MkdirAll
var GlobalGormDatabase *GormDatabase

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, logLevel string) (*GormDatabase, error) {
	createDirectoryIfSqlite(dialect, connection)
	_connection := fmt.Sprintf("%s?_foreign_keys=on", connection)
	db, err := gorm.Open(sqlite.Open(_connection), &gorm.Config{
		Logger: logger.New().SetLogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("invalid database")
	}
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = migration.AutoMigrate(db); err != nil {
		panic("failed to AutoMigrate")
	}

	user_, _ := user.GetUser()
	if user_ == nil {
		_, _ = user.CreateUser(&user.User{Username: username, Password: password})
	}

	busService := eventbus.NewService(eventbus.GetBus())
	rubixRegistry := rubixregistry.New()
	GlobalGormDatabase = &GormDatabase{DB: db, Bus: busService, RubixRegistry: rubixRegistry}
	return &GormDatabase{DB: db, Bus: busService, RubixRegistry: rubixRegistry}, nil
}

func createDirectoryIfSqlite(dialect, connection string) {
	if dialect == "sqlite3" {
		if _, err := os.Stat(filepath.Dir(connection)); os.IsNotExist(err) {
			if err := mkdirAll(filepath.Dir(connection), 0777); err != nil {
				panic(err)
			}
		}
	}
}

// GormDatabase is a wrapper for the gorm framework.
type GormDatabase struct {
	DB            *gorm.DB
	Store         cachestore.Handler
	Bus           eventbus.BusService
	PluginManager *plugin.Manager
	Modules       map[string]shared.Module
	RubixRegistry *rubixregistry.RubixRegistry
}

// Close closes the gorm database connection.
func (d *GormDatabase) Close() {
	sqlDB, err := d.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
