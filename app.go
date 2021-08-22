package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/mode"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/router"
	"github.com/NubeDev/flow-framework/runner"
)

var (
	Version = "unknown"
	// Commit the git commit hash of this version.
	Commit = "unknown"
	// BuildDate the date on which this binary was build.
	BuildDate = "unknown"
	// Mode the build mode.
	Mode = mode.Dev
)



func main() {
	vInfo := &model.VersionInfo{Version: Version, Commit: Commit, BuildDate: BuildDate}
	mode.Set(Mode)

	log.Println("Starting version", vInfo.Version+"@"+BuildDate)
	rand.Seed(time.Now().UnixNano())

	p := flag.String("config", "./config.yml", "config file path")
	flag.Parse()
	path := *p
	conf := config.Get(path)

	if conf.PluginsDir != "" {
		if err := os.MkdirAll(conf.PluginsDir, 0755); err != nil {
			panic(err)
		}
	}
	if err := os.MkdirAll(conf.UploadedImagesDir, 0755); err != nil {
		panic(err)
	}

	db, err := database.New(conf.Database.Dialect, conf.Database.Connection, conf.DefaultUser.Name, conf.DefaultUser.Pass, conf.PassStrength, true)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	engine, closeable := router.Create(db, vInfo, conf)
	defer closeable()

	runner.Run(engine, conf)
}
