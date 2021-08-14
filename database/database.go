package database

import (
	"fmt"
	"github.com/NubeDev/plug-framework/auth/password"
	"github.com/NubeDev/plug-framework/model"
	"gorm.io/driver/sqlite"
	"os"
	"path/filepath"
	_ "time"

	_ "github.com/NubeDev/plug-framework/auth/password"
	_ "github.com/NubeDev/plug-framework/model"
	"gorm.io/gorm"

	// enable the mysql dialect.
	_ "gorm.io/driver/mysql"

	// enable the postgres dialect.
	_ "gorm.io/driver/postgres"

	// enable the sqlite3 dialect.
	_ "gorm.io/driver/sqlite"
)

var mkdirAll = os.MkdirAll

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, defaultUser, defaultPass string, strength int, createDefaultUserIfNotExist bool) (*GormDatabase, error) {
	createDirectoryIfSqlite(dialect, connection)

	fmt.Println(connection)
	//path := fmt.Sprintf("%s?_foreign_keys=on", connection)
	//db, err := gorm.Open(dialect, connection)
	db, err := gorm.Open(sqlite.Open(connection), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var User []model.User
	var Application []model.Application
	var Message []model.Message
	var Client []model.Client
	var PluginConf []model.PluginConf
	var models = []interface{}{
		&User,
		&Application,
		&Message,
		&Client,
		&PluginConf,
	}

	//if err = db.AutoMigrate(new(model.User), new(model.VersionInfo),  new(model.Application), new(model.Message), new(model.Client), new(model.PluginConf)).Error; err != nil {
	//	return nil, err
	//}

	for _, v := range models {
		err = db.AutoMigrate(v)
		if err != nil {
			fmt.Println(err)
			panic("failed to AutoMigrate")
			//fmt.Println("db migrate issue")
		}
	}

	// We normally don't need that much connections, so we limit them. F.ex. mysql complains about
	// "too many connections", while load testing Gotify.
	//db.DB().SetMaxOpenConns(10)
	//
	//if dialect == "sqlite3" {
	//	// We use the database connection inside the handlers from the http
	//	// framework, therefore concurrent access occurs. Sqlite cannot handle
	//	// concurrent writes, so we limit sqlite to one connection.
	//	// see https://github.com/mattn/go-sqlite3/issues/274
	//	db.DB().SetMaxOpenConns(1)
	//}
	//
	//if dialect == "mysql" {
	//	// Mysql has a setting called wait_timeout, which defines the duration
	//	// after which a connection may not be used anymore.
	//	// The default for this setting on mariadb is 10 minutes.
	//	// See https://github.com/docker-library/mariadb/issues/113
	//	db.DB().SetConnMaxLifetime(9 * time.Minute)
	//}
	//
	//if err = db.AutoMigrate(new(model.User), new(model.VersionInfo),  new(model.Application), new(model.Message), new(model.Client), new(model.PluginConf)).Error; err != nil {
	//	return nil, err
	//}

	//if err = prepareBlobColumn(dialect, db); err != nil {
	//	return nil, err
	//}

	var userCount int64 = 0
	db.Find(new(model.User)).Count(&userCount)
	if createDefaultUserIfNotExist && userCount == 0 {
		db.Create(&model.User{Name: defaultUser, Pass: password.CreatePassword(defaultPass, strength), Admin: true})
	}

	return &GormDatabase{DB: db}, nil
}

func prepareBlobColumn(dialect string, db *gorm.DB) error {
	//blobType := ""
	//switch dialect {
	//case "mysql":
	//	blobType = "longblob"
	//case "postgres":
	//	blobType = "bytea"
	//}
	//if blobType != "" {
	//	for _, target := range []struct {
	//		Table  interface{}
	//		Column string
	//	}{
	//		{model.Message{}, "extras"},
	//		{model.PluginConf{}, "config"},
	//		{model.PluginConf{}, "storage"},
	//	} {
	//		//if err := db.Model(target.Table).ModifyColumn(target.Column, blobType).Error; err != nil {
	//		//	return err
	//		//}
	//	}
	//}
	return nil
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
	d.Close()
	//d.DB.Close()
}
