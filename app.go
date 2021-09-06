package main

import (
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/logger"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/router"
	"github.com/NubeDev/flow-framework/runner"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

var (
	Version   = "<version>"
	Commit    = "<commit>"
	BuildDate = "<build_date>"
)

func main() {
	conf := config.CreateApp()
	logger.SetLogger(conf.LogLevel)
	logger.SetGinMode(conf.LogLevel)

	vInfo := &model.VersionInfo{Version: Version, Commit: Commit, BuildDate: BuildDate}
	log.Info("Info Starting version:", vInfo.Version+"-"+vInfo.Commit+"@"+vInfo.BuildDate)

	eventbus.InitBus()
	database.DataBus()
	if err := os.MkdirAll(conf.GetAbsPluginDir(), 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(conf.GetAbsUploadedImagesDir(), 0755); err != nil {
		panic(err)
	}
	connection := path.Join(conf.GetAbsDataDir(), conf.Database.Connection)
	db, err := database.New(conf.Database.Dialect, connection, conf.DefaultUser.Name, conf.DefaultUser.Pass,
		conf.PassStrength, conf.Database.LogLevel, true)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	engine, closeable := router.Create(db, vInfo, conf)

	defer closeable()
	runner.Run(engine, conf)

}
