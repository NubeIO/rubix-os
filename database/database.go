package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/auth/password"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/src/cachestore"
	"github.com/NubeDev/flow-framework/utils"
	"os"
	"path/filepath"

	"github.com/NubeDev/flow-framework/logger"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var mkdirAll = os.MkdirAll
var gormDatabase *GormDatabase

//var GetDatabaseBus eventbus.BusService

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, defaultUser, defaultPass string, strength int, logLevel string,
	createDefaultUserIfNotExist bool, production bool) (*GormDatabase, error) {
	createDirectoryIfSqlite(dialect, connection)
	_connection := fmt.Sprintf("%s?_foreign_keys=on", connection)
	db, err := gorm.Open(sqlite.Open(_connection), &gorm.Config{
		Logger: logger.New().SetLogMode(logLevel),
	})
	if err != nil {
		panic("failed to connect database")
	}
	var localStorageFlowNetwork []model.LocalStorageFlowNetwork
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
	var flowNetworkClone []model.FlowNetworkClone
	var job []model.Job
	var stream []model.Stream
	var streamClone []model.StreamClone
	var commandGroup []model.CommandGroup
	var producer []model.Producer
	var consumer []model.Consumer
	var writer []model.Writer
	var writerClone []model.WriterClone
	var integration []model.Integration
	var mqttConnection []model.MqttConnection
	var serialConnection []model.SerialConnection
	var ipConnection []model.IpConnection
	var schedule []model.Schedule
	var tags []model.Tag
	var blocks []model.Block
	var connections []model.Connection
	var blockRoutes []model.BlockStaticRoute
	var blockRouteValueNumbers []model.BlockRouteValueNumber
	var blockRouteValueString []model.BlockRouteValueString
	var blockRouteValueBool []model.BlockRouteValueBool
	var sourceParams []model.SourceParameter
	var links []model.Link
	var history []model.History
	var historyLog []model.HistoryLog
	var models = []interface{}{
		&localStorageFlowNetwork,
		&alerts,
		&user,
		&application,
		&message,
		&client,
		&pluginConf,
		&network,
		&device,
		&point,
		&flowNetwork,
		&flowNetworkClone,
		&priority,
		&producerHistory,
		&consumerHistory,
		&job,
		&stream,
		&streamClone,
		&commandGroup,
		&producer,
		&consumer,
		&writer,
		&writerClone,
		&integration,
		&mqttConnection,
		&serialConnection,
		&ipConnection,
		&schedule,
		&tags,
		&blocks,
		&connections,
		&blockRoutes,
		&blockRouteValueNumbers,
		&blockRouteValueString,
		&blockRouteValueBool,
		&sourceParams,
		&links,
		&history,
		&historyLog,
	}

	for _, v := range models {
		err = db.AutoMigrate(v)
		if err != nil {
			panic("failed to AutoMigrate")
		}
	}

	var userCount int64 = 0
	db.Find(new(model.User)).Count(&userCount)
	if createDefaultUserIfNotExist && userCount == 0 {
		db.Create(&model.User{Name: defaultUser, Pass: password.CreatePassword(defaultPass, strength), Admin: true})
		//if !production { //make a fake token for dev
		//	c := new(model.Client)
		//	c.Token = "fakeToken123"
		//	c.UserID = 1
		//	c.Name = "admin"
		//	db.Create(c)
		//}
	}
	var lsFlowNetworkCount int64 = 0
	lsFlowNetwork := model.LocalStorageFlowNetwork{
		FlowIP:    "0.0.0.0",
		FlowPort:  1660,
		FlowHTTPS: utils.NewFalse(),
	}
	db.Find(&model.LocalStorageFlowNetwork{}).Count(&lsFlowNetworkCount)
	if lsFlowNetworkCount == 0 {
		db.Create(&lsFlowNetwork)
	}
	busService := eventbus.NewService(eventbus.GetBus())
	gormDatabase = &GormDatabase{DB: db, Bus: busService}
	return &GormDatabase{DB: db, Bus: busService}, nil
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
}

// Close closes the gorm database connection.
func (d *GormDatabase) Close() {
	fmt.Println("FIX THIS CLOSE DB dont work after upgrade of gorm")
	//d.Close() //TODO this is broken it calls itself
}
