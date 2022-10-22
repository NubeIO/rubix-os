package config

import (
	"flag"
	"path"
	"time"

	"github.com/NubeIO/configor"
)

type Configuration struct {
	Server struct {
		KeepAlivePeriodSeconds int
		ListenAddr             string `default:"0.0.0.0"`
		Port                   int
		RSPort                 int `default:"1616"`
		ResponseHeaders        map[string]string
		Stream                 struct {
			PingPeriodSeconds int `default:"45"`
			AllowedOrigins    []string
		}
		Cors struct {
			AllowOrigins []string
			AllowMethods []string
			AllowHeaders []string
		}
	}
	Database struct {
		Dialect    string `default:"sqlite3"`
		Connection string `default:"data.db"`
		LogLevel   string `default:"WARN"`
	}
	PassStrength int `default:"10"`
	LogLevel     string
	Location     struct {
		GlobalDir string `default:"./"`
		ConfigDir string `default:"config"`
		DataDir   string `default:"data"`
		Data      struct {
			PluginsDir        string `default:"plugins"`
			UploadedImagesDir string `default:"images"`
		}
		DeviceInfoFile string `default:"/data/rubix-registry/device_info.json"`
	}
	Prod            bool  `default:"false"`
	Auth            *bool `default:"true"`
	ProducerHistory struct {
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
	SecretKey string
	MQTT      struct {
		Enable         *bool         `default:"true"`
		Address        string        `default:"localhost"`
		Port           int           `default:"1883"`
		Username       string        `default:""`
		Password       string        `default:""`
		SetKeepAlive   time.Duration `default:"0s"`
		SetPingTimeout time.Duration `default:"0s"`
		AutoReconnect  bool          `default:"true"`
		QOS            int           `default:"1"`
		Retain         *bool         `default:"true"`
	}
}

var config *Configuration = nil

func Get() *Configuration {
	return config
}

func CreateApp() *Configuration {
	config = new(Configuration)
	config = config.Parse()
	err := configor.New(&configor.Config{EnvironmentPrefix: "FLOW"}).Load(config, path.Join(config.GetAbsConfigDir(), "config.yml"))
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

func (conf *Configuration) GetAbsPluginDir() string {
	return path.Join(conf.GetAbsDataDir(), conf.Location.Data.PluginsDir)
}

func (conf *Configuration) GetAbsUploadedImagesDir() string {
	return path.Join(conf.GetAbsDataDir(), conf.Location.Data.UploadedImagesDir)
}
