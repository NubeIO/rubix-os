package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/NubeDev/plug-framework/config"
	"github.com/NubeDev/plug-framework/database"
	"github.com/NubeDev/plug-framework/mode"
	"github.com/NubeDev/plug-framework/model"
	"github.com/NubeDev/plug-framework/router"
	"github.com/NubeDev/plug-framework/runner"
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

	fmt.Println("Starting version", vInfo.Version+"@"+BuildDate)
	rand.Seed(time.Now().UnixNano())
	conf := config.Get()

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
