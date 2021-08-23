package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/auth/password"
	"github.com/NubeDev/flow-framework/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
)

var mkdirAll = os.MkdirAll

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, defaultUser, defaultPass string, strength int, createDefaultUserIfNotExist bool) (*GormDatabase, error) {
	createDirectoryIfSqlite("sqlite3", connection)
	_connection := fmt.Sprintf("%s?_foreign_keys=on", connection)
	db, err := gorm.Open(sqlite.Open(_connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	var user []model.User
	var application []model.Application
	var message []model.Message
	var client []model.Client
	var pluginConf []model.PluginConf
	var network []model.Network
	var device []model.Device
	var point []model.Point
	//var pointStore []model.PointStore
	//var priorityArrayModel []model.PriorityArrayModel
	var job []model.Job
	var gateway []model.Gateway
	var subscriber []model.Subscriber
	var subscription []model.Subscription
	var models = []interface{}{
		&user,
		&application,
		&message,
		&client,
		&pluginConf,
		&network,
		&device,
		&point,
		//&pointStore,
		//&priorityArrayModel,
		&job,
		&gateway,
		&subscriber,
		&subscription,
	}

	for _, v := range models {
		err = db.AutoMigrate(v)
		if err != nil {
			fmt.Println(err)
			panic("failed to AutoMigrate")
		}
	}

	var userCount int64 = 0
	db.Find(new(model.User)).Count(&userCount)
	if createDefaultUserIfNotExist && userCount == 0 {
		db.Create(&model.User{Name: defaultUser, Pass: password.CreatePassword(defaultPass, strength), Admin: true})
	}

	return &GormDatabase{DB: db}, nil
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
	DB *gorm.DB
}


// Close closes the gorm database connection.
func (d *GormDatabase) Close() {
	fmt.Println("FIX THIS CLOSE DB dont work after upgrade of gorm")
	//d.Close() //TODO this is broken it calls itself
}
