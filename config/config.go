package config

import (
	"flag"
	"github.com/NubeDev/configor"
	"path"
)

type Configuration struct {
	Server struct {
		KeepAlivePeriodSeconds int
		ListenAddr             string `default:"0.0.0.0"`
		Port                   int
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
	DefaultUser struct {
		Name string `default:"admin"`
		Pass string `default:"admin"`
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
	}
	Prod bool `default:"false"`
	PointBuilder
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
	flag.Parse()
	conf.Server.Port = *port
	conf.Location.GlobalDir = *globalDir
	conf.Location.DataDir = *dataDir
	conf.Location.ConfigDir = *configDir
	conf.Prod = *prod
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

type PointBuilder struct {
	TotalWaterConsumption struct {
		ThingClass   string   `yaml:"thingClass"`
		ThingType    string   `yaml:"thingType"`
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"TotalWaterConsumption"`
	WaterTarget struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"WaterTarget"`
	TotalGasConsumption struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"TotalGasConsumption"`
	GasTarget struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"GasTarget"`
	TotalActiveEnergy struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"TotalActiveEnergy"`
	TotalActivePower struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"TotalActivePower"`
	EnergyTarget struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"EnergyTarget"`
	ZoneRoomTemperature struct {
		ThingClass   string   `yaml:"thingClass"`
		ThingType    string   `yaml:"thingType"`
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ZoneRoomTemperature"`
	ZoneRoomTemperatureSetPoint struct {
		ThingClass   string   `yaml:"thingClass"`
		ThingType    string   `yaml:"thingType"`
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ZoneRoomTemperatureSetPoint"`
	SupplyAirTemperature struct {
		ThingClass   string   `yaml:"thingClass"`
		ThingType    string   `yaml:"thingType"`
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirTemperature"`
	ReturnAirTemperature struct {
		ThingClass   string   `yaml:"thingClass"`
		ThingType    string   `yaml:"thingType"`
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirTemperature"`
	ZoneRoomRelativeHumidity struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ZoneRoomRelativeHumidity"`
	Damper struct {
		Kind      string   `yaml:"kind"`
		Unit      string   `yaml:"unit"`
		EquipType string   `yaml:"equipType"`
		Tags      []string `yaml:"tags"`
	} `yaml:"Damper"`
	Airflow struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"Airflow"`
	AirflowSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AirflowSetPoint"`
	MinimumAirflow struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"MinimumAirflow"`
	ElectricDuctHeaterStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ElectricDuctHeaterStatus"`
	ElectricDuctHeaterStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ElectricDuctHeaterStartStop"`
	TerminalLoad struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"TerminalLoad"`
	MaximumAirflow struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"MaximumAirflow"`
	GeneralExhaustAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"GeneralExhaustAirFanStartStop"`
	GeneralExhaustAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"GeneralExhaustAirFanStatus"`
	GeneralExhaustAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"GeneralExhaustAirFanFault"`
	GeneralExhaustAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"GeneralExhaustAirconstPressure"`
	ExhaustAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ExhaustAirFanStartStop"`
	ExhaustAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ExhaustAirFanStatus"`
	ExhaustAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ExhaustAirFanFault"`
	ExhaustAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ExhaustAirconstPressure"`
	KitchenExhaustAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"KitchenExhaustAirFanStartStop"`
	KitchenExhaustAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"KitchenExhaustAirFanStatus"`
	KitchenExhaustAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"KitchenExhaustAirFanFault"`
	KitchenExhaustAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"KitchenExhaustAirconstPressure"`
	ToiletExhaustAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ToiletExhaustAirFanStartStop"`
	ToiletExhaustAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ToiletExhaustAirFanStatus"`
	ToiletExhaustAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ToiletExhaustAirFanFault"`
	ToiletExhaustAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ToiletExhaustAirconstPressure"`
	CarParkExhaustAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkExhaustAirFanStartStop"`
	CarParkExhaustAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkExhaustAirFanStatus"`
	CarParkExhaustAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkExhaustAirFanFault"`
	CarParkExhaustAirFanSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkExhaustAirFanSpeed"`
	CarParkExhaustAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkExhaustAirconstPressure"`
	CarParkExhaustAirconstPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkExhaustAirconstPressureSetPoint"`
	CarParkExhaustCarbonMonoxide struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkExhaustCarbonMonoxide"`
	CarParkExhaustCarbonMonoxideSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkExhaustCarbonMonoxideSetPoint"`
	CarParkSupplyAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkSupplyAirFanStartStop"`
	CarParkSupplyAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkSupplyAirFanStatus"`
	CarParkSupplyAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkSupplyAirFanFault"`
	CarParkSupplyAirFanSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkSupplyAirFanSpeed"`
	CarParkSupplyAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkSupplyAirconstPressure"`
	CarParkSupplyAirconstPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CarParkSupplyAirconstPressureSetPoint"`
	AHUReturnAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUReturnAirFanStartStop"`
	AHUReturnAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUReturnAirFanStatus"`
	AHUReturnAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUReturnAirFanFault"`
	AHUReturnAirFanSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUReturnAirFanSpeed"`
	AHUReturnAirTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUReturnAirTemperature"`
	AHUReturnAirHumidity struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUReturnAirHumidity"`
	AHUReturnAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUReturnAirconstPressure"`
	AHUReturnAirCO2 struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUReturnAirCO2"`
	AHUReturnAirconstPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUReturnAirconstPressureSetPoint"`
	ReturnAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirFanStartStop"`
	ReturnAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirFanStatus"`
	ReturnAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirFanFault"`
	ReturnAirFanSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirFanSpeed"`
	ReturnAirHumidity struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirHumidity"`
	ReturnAirFlow struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirFlow"`
	ReturnAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirconstPressure"`
	ReturnAirEnthalpy struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirEnthalpy"`
	ReturnAirSuctionPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirSuctionPressure"`
	ReturnAirCO2 struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirCO2"`
	ReturnAirconstPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirconstPressureSetPoint"`
	SupplyAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirFanStartStop"`
	SupplyAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirFanStatus"`
	SupplyAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirFanFault"`
	SupplyAirFanSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirFanSpeed"`
	SupplyAirFlow struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirFlow"`
	SupplyAirHumidity struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirHumidity"`
	SupplyAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirconstPressure"`
	SupplyAirconstPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirconstPressureSetPoint"`
	OutsideAirFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirFanStartStop"`
	OutsideAirFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirFanStatus"`
	OutsideAirFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirFanFault"`
	OutsideAirFanSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirFanSpeed"`
	OutsideAirTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirTemperature"`
	OutsideAirHumidity struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirHumidity"`
	OutsideAirFlow struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirFlow"`
	OutsideAirconstPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirconstPressure"`
	OutsideAirconstPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirconstPressureSetPoint"`
	PACFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"PACFanStartStop"`
	PACFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"PACFanStatus"`
	PACFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"PACFanFault"`
	ReturnRoomAirTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnRoomAirTemperature"`
	ReturnRoomAirTemperatureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnRoomAirTemperatureSetPoint"`
	ReturnRoomAirHumidity struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnRoomAirHumidity"`
	CompressorStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CompressorStartStop"`
	CondenserWaterCall struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterCall"`
	AirOffCoilTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AirOffCoilTemperature"`
	OutsideAirDamperControl struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirDamperControl"`
	FilterDifferentialPressureStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FilterDifferentialPressureStatus"`
	ReverseValveControl struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReverseValveControl"`
	HeatingCoilValveControl struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HeatingCoilValveControl"`
	FTUFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FTUFanStartStop"`
	DamperFeedback struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"DamperFeedback"`
	FCUFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FCUFanStartStop"`
	FCUFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FCUFanStatus"`
	FCUFanSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FCUFanSpeed"`
	FCUFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FCUFanFault"`
	SupplyAirTemperatureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirTemperatureSetPoint"`
	CoolingCoilValveControl struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CoolingCoilValveControl"`
	AHUFanStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUFanStartStop"`
	AHUFanStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUFanStatus"`
	AHUFanFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUFanFault"`
	Speed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"Speed"`
	RoomAirTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"RoomAirTemperature"`
	MixAirTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"MixAirTemperature"`
	ReturnAirTemperatureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirTemperatureSetPoint"`
	RoomAirTemperatureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"RoomAirTemperatureSetPoint"`
	RoomAirHumidity struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"RoomAirHumidity"`
	OutsideAirRelativeHumidity struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirRelativeHumidity"`
	MinimumOutsideAirDamperControl struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"MinimumOutsideAirDamperControl"`
	MinimumOutsideAirDamperFeedback struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"MinimumOutsideAirDamperFeedback"`
	ExhaustAirDamper struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ExhaustAirDamper"`
	RACO2 struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"RACO2"`
	DewPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"DewPoint"`
	DewPointSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"DewPointSetPoint"`
	AHUFanSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHUFanSpeed"`
	OutsideReturnRoomAirTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideReturnRoomAirTemperature"`
	SupplyReturnRoomAirHumidity struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyReturnRoomAirHumidity"`
	ReturnAirDamperControl struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirDamperControl"`
	FilterDifferentialPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FilterDifferentialPressure"`
	FilterDifferentialPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FilterDifferentialPressureSetPoint"`
	FilterStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FilterStatus"`
	FloorOccupancy struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"FloorOccupancy"`
	AfterHoursActiveTime struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AfterHoursActiveTime"`
	ZoneOccupancy struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ZoneOccupancy"`
	ZoneAfterHoursActiveTime struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ZoneAfterHoursActiveTime"`
	AfterHoursElapsedTime struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AfterHoursElapsedTime"`
	AHEnable struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"AHEnable"`
	VAVEnable struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"VAVEnable"`
	ReturnAirDamper struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ReturnAirDamper"`
	SupplyAirDamper struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SupplyAirDamper"`
	IsolationAirDamper struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"IsolationAirDamper"`
	ZoneTemperatureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ZoneTemperatureSetPoint"`
	CommonHotWaterSupplyTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonHotWaterSupplyTemperature"`
	CommonHotWaterReturnTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonHotWaterReturnTemperature"`
	HotWaterPumpStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HotWaterPumpStartStop"`
	HotWaterPumpStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HotWaterPumpStatus"`
	HotWaterPumpFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HotWaterPumpFault"`
	HotWaterPumpDifferentialPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HotWaterPumpDifferentialPressure"`
	HotWaterPumpDifferentialPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HotWaterPumpDifferentialPressureSetPoint"`
	HotWaterPumpSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HotWaterPumpSpeed"`
	BoilerStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"BoilerStartStop"`
	BoilerStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"BoilerStatus"`
	BoilerFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"BoilerFault"`
	HotWaterSupplyTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HotWaterSupplyTemperature"`
	HotWaterSupplyTemperatureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HotWaterSupplyTemperatureSetPoint"`
	HotWaterReturnTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"HotWaterReturnTemperature"`
	CommonCondenserWaterSupplyTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonCondenserWaterSupplyTemperature"`
	CommonCondenserPressureDiff struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonCondenserPressureDiff"`
	CommonBypassValve struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonBypassValve"`
	CommonCondenserWaterReturnTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonCondenserWaterReturnTemperature"`
	CondenserWaterPumpStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterPumpStartStop"`
	CondenserWaterPumpStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterPumpStatus"`
	CondenserWaterPumpFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterPumpFault"`
	CondenserWaterPumpDifferentialPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterPumpDifferentialPressure"`
	CondenserWaterPumpDifferentialPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterPumpDifferentialPressureSetPoint"`
	CondenserWaterFlow struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterFlow"`
	CondenserWaterPumpSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterPumpSpeed"`
	CoolingTowerStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CoolingTowerStartStop"`
	CoolingTowerStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CoolingTowerStatus"`
	CoolingTowerFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CoolingTowerFault"`
	CoolingTowerFanSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CoolingTowerFanSpeed"`
	CondenserWaterSupplyTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterSupplyTemperature"`
	CondenserWaterSupplyTemperatureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterSupplyTemperatureSetPoint"`
	CondenserWaterReturnTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CondenserWaterReturnTemperature"`
	PrimarySupplyTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"PrimarySupplyTemperature"`
	PrimaryReturnTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"PrimaryReturnTemperature"`
	PrimaryLoopValve struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"PrimaryLoopValve"`
	SecondarySupplyTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SecondarySupplyTemperature"`
	SecondaryReturnTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SecondaryReturnTemperature"`
	SecondaryLoopValve struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"SecondaryLoopValve"`
	OutsideAirEnthalpy struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"OutsideAirEnthalpy"`
	CommonChilledWaterSupplyTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonChilledWaterSupplyTemperature"`
	CommonChilledWaterReturnTemperature struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonChilledWaterReturnTemperature"`
	CommonChilledWaterDifferentialPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonChilledWaterDifferentialPressure"`
	CommonChilledWaterBypassValve struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonChilledWaterBypassValve"`
	CommonChilledWaterDifferentialPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"CommonChilledWaterDifferentialPressureSetPoint"`
	ChilledWaterPumpStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChilledWaterPumpStartStop"`
	ChilledWaterPumpStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChilledWaterPumpStatus"`
	ChilledWaterPumpFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChilledWaterPumpFault"`
	ChilledWaterPumpDifferentialPressure struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChilledWaterPumpDifferentialPressure"`
	ChilledWaterPumpDifferentialPressureSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChilledWaterPumpDifferentialPressureSetPoint"`
	ChilledWaterFlow struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChilledWaterFlow"`
	ChilledWaterFlowSetPoint struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChilledWaterFlowSetPoint"`
	ChilledWaterPumpSpeed struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChilledWaterPumpSpeed"`
	ChillerStartStop struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChillerStartStop"`
	ChillerStatus struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChillerStatus"`
	ChillerFault struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"ChillerFault"`
	Load struct {
		Kind         string   `yaml:"kind"`
		Unit         string   `yaml:"unit"`
		UnitImperial string   `yaml:"unitImperial"`
		EquipType    string   `yaml:"equipType"`
		Tags         []string `yaml:"tags"`
	} `yaml:"Load"`
	ChilledWaterSupplyTemperature struct {
		Kind      string   `yaml:"kind"`
		Unit      string   `yaml:"unit"`
		EquipType string   `yaml:"equipType"`
		Tags      []string `yaml:"tags"`
	} `yaml:"ChilledWaterSupplyTemperature"`
	ChilledWaterSupplyTemperatureSetPoint struct {
		Kind      string   `yaml:"kind"`
		Unit      string   `yaml:"unit"`
		EquipType string   `yaml:"equipType"`
		Tags      []string `yaml:"tags"`
	} `yaml:"ChilledWaterSupplyTemperatureSetPoint"`
	ChilledWaterReturnTemperature struct {
		Kind      string   `yaml:"kind"`
		Unit      string   `yaml:"unit"`
		EquipType string   `yaml:"equipType"`
		Tags      []string `yaml:"tags"`
	} `yaml:"ChilledWaterReturnTemperature"`
}
