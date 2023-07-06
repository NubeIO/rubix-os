package config

import (
	"flag"
	"github.com/NubeIO/configor"
	"path"
	"time"
)

type Configuration struct {
	Server struct {
		KeepAlivePeriodSeconds int
		ListenAddr             string `default:"0.0.0.0"`
		Port                   int
		ResponseHeaders        map[string]string
	}
	Database struct {
		Dialect    string `default:"sqlite3"`
		Connection string `default:"data.db"`
		LogLevel   string `default:"WARN"`
	}
	LogLevel string
	Location struct {
		GlobalDir string `default:"./"`
		ConfigDir string `default:"config"`
		DataDir   string `default:"data"`
		Data      struct {
			PluginsDir        string `default:"plugins"`
			ModulesDir        string `default:"modules"`
			UploadedImagesDir string `default:"images"`
		}
	} // leave it as default; don't include in config.eg.yml
	Prod         bool  `default:"false"` // set from commandline; don't include in config.eg.yml
	Auth         *bool `default:"true"`  // set from commandline; don't include in config.eg.yml
	PointHistory struct {
		Enable  *bool `default:"true"`
		Cleaner struct {
			Enable              *bool `default:"true"`
			Frequency           int   `default:"600"`
			DataPersistingHours int   `default:"24"`
		}
		IntervalHistoryCreator struct {
			Enable    *bool `default:"true"`
			Frequency int   `default:"10"`
		}
	}
	MQTT struct {
		Enable                *bool  `default:"true"`
		Address               string `default:"localhost"`
		Port                  int    `default:"1883"`
		Username              string `default:""`
		Password              string `default:""`
		AutoReconnect         *bool  `default:"true"`
		ConnectRetry          *bool  `default:"true"`
		QOS                   int    `default:"1"`
		Retain                *bool  `default:"true"`
		GlobalBroadcast       *bool  `default:"false"` // if set to true will include the plat details in the topic
		PublishPointCOV       *bool  `default:"true"`
		PublishPointList      *bool  `default:"false"`
		PointWriteListener    *bool  `default:"true"`
		PublishScheduleCOV    *bool  `default:"true"`
		PublishScheduleList   *bool  `default:"false"`
		ScheduleWriteListener *bool  `default:"true"`
	}
	Notification struct {
		Enable         *bool         `default:"false"`
		Frequency      time.Duration `default:"1m"`
		ResendDuration time.Duration `default:"1h"`
	}
}

var config *Configuration = nil

func Get() *Configuration {
	return config
}

func CreateApp() *Configuration {
	config = new(Configuration)
	config = config.Parse()
	err := configor.New(&configor.Config{EnvironmentPrefix: "ROS"}).Load(config, path.Join(config.GetAbsConfigDir(), "config.yml"))
	if err != nil {
		panic(err)
	}
	return config
}

func (conf *Configuration) Parse() *Configuration {
	port := flag.Int("p", 1660, "Port")
	globalDir := flag.String("g", "./", "Global Directory")
	dataDir := flag.String("d", "data", "Data Directory")
	configDir := flag.String("c", "config", "Config Directory")
	prod := flag.Bool("prod", false, "Deployment Mode")
	auth := flag.Bool("auth", true, "enable auth")
	flag.Parse()
	conf.Server.Port = *port
	conf.Location.GlobalDir = *globalDir
	conf.Location.DataDir = *dataDir
	conf.Location.ConfigDir = *configDir
	conf.Prod = *prod
	conf.Auth = auth
	return conf
}

func (conf *Configuration) GetAbsDataDir() string {
	return path.Join(conf.Location.GlobalDir, conf.Location.DataDir)
}

func (conf *Configuration) GetAbsConfigDir() string {
	return path.Join(conf.Location.GlobalDir, conf.Location.ConfigDir)
}

func (conf *Configuration) GetAbsPluginsDir() string {
	return path.Join(conf.GetAbsDataDir(), conf.Location.Data.PluginsDir)
}

func (conf *Configuration) GetAbsModulesDir() string {
	return path.Join(conf.GetAbsDataDir(), conf.Location.Data.ModulesDir)
}

func (conf *Configuration) GetAbsUploadedImagesDir() string {
	return path.Join(conf.GetAbsDataDir(), conf.Location.Data.UploadedImagesDir)
}

func (conf *Configuration) GetAbsTempDir() string {
	return "/tmp"
}

func (conf *Configuration) GetSnapshotDir() string {
	return "/data"
}

func (conf *Configuration) GetAbsSnapShotDir() string {
	return path.Join(conf.GetAbsDataDir(), "snapshots")
}
