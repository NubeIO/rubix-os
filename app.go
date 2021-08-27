package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/router"
	"github.com/NubeDev/flow-framework/runner"
	"os"
	"path"
)

var (
	Version   = "<version>"
	Commit    = "<commit>"
	BuildDate = "<build_date>"
)



func main() {
	vInfo := &model.VersionInfo{Version: Version, Commit: Commit, BuildDate: BuildDate}
	fmt.Println("Starting version:", vInfo.Version+"-"+vInfo.Commit+"@"+vInfo.BuildDate)
	// Start Event Bus
	eventbus.InitBus()
	database.DataBus()
	//config
	conf := config.CreateApp()


	if err := os.MkdirAll(conf.GetAbsPluginDir(), 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(conf.GetAbsUploadedImagesDir(), 0755); err != nil {
		panic(err)
	}
	//db
	connection := path.Join(conf.GetAbsDataDir(), conf.Database.Connection)
	db, err := database.New(conf.Database.Dialect, connection, conf.DefaultUser.Name, conf.DefaultUser.Pass, conf.PassStrength, true)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//router
	engine, closeable := router.Create(db, vInfo, conf)
	defer closeable()
	//run
	runner.Run(engine, conf)
}
