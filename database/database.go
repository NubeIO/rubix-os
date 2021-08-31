package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/auth/password"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
)

var mkdirAll = os.MkdirAll

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, defaultUser, defaultPass string, strength int, createDefaultUserIfNotExist bool) (*GormDatabase, error) {
	createDirectoryIfSqlite(dialect, connection)
	_connection := fmt.Sprintf("%s?_foreign_keys=on", connection)
	db, err := gorm.Open(sqlite.Open(_connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	var alerts []model.Alert
	var user []model.User
	var application []model.Application
	var message []model.Message
	var client []model.Client
	var pluginConf []model.PluginConf
	var network []model.Network
	var device []model.Device
	var point []model.Point
	var priority []model.Priority
	var producerHistory []model.ProducerHistory
	var consumerHistory []model.ConsumerHistory
	var flowNetwork []model.FlowNetwork
	var rubixPlat []model.RubixPlat
	var job []model.Job
	var stream []model.Stream
	var streamList []model.StreamList
	var commandGroup []model.CommandGroup
	var producer []model.Producer
	var consumer []model.Consumer
	var writer []model.Writer
	var writerCopy []model.WriterClone

	var models = []interface{}{
		&alerts,
		&user,
		&application,
		&message,
		&client,
		&pluginConf,
		&network,
		&device,
		&point,
		&rubixPlat,
		&flowNetwork,
		&priority,
		&producerHistory,
		&consumerHistory,
		&job,
		&stream,
		&streamList,
		&commandGroup,
		&producer,
		&consumer,
		&writer,
		&writerCopy,
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
	var platCount int64 = 0
	rp := new(model.RubixPlat)
	db.Find(rp).Count(&platCount)
	if createDefaultUserIfNotExist && platCount == 0 {
		rp.GlobalUuid = utils.MakeTopicUUID(model.CommonNaming.Rubix)
		db.Create(&rp)
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
