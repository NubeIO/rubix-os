package main

type Config struct {
	Postgres           Postgres           `yaml:"postgres"`
	CloudServerDetails CloudServerDetails `yaml:"cloudServerDetails"`
	Job                Job                `yaml:"job"`
	LogLevel           string             `yaml:"log_level"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
	Port     int    `yaml:"port"`
	SslMode  string `yaml:"ssl_mode"`
}

type Job struct {
	Frequency        string `yaml:"frequency"`
	SyncPointsWithDB bool   `yaml:"sync_points_with_db"`
}

type CloudServerDetails struct {
	CloudHostUUID string `yaml:"cloudHostUUID"`
}

func (inst *Instance) DefaultConfig() interface{} {
	postgres := Postgres{
		Host:     "localhost",
		Port:     5432,
		DbName:   "db_ff",
		User:     "postgres",
		Password: "password",
		SslMode:  "disable",
	}
	cloudServerDetails := CloudServerDetails{
		CloudHostUUID: "hos_afccc787e32f411e",
	}
	job := Job{
		Frequency:        "1m",
		SyncPointsWithDB: false,
	}

	return &Config{
		Postgres:           postgres,
		CloudServerDetails: cloudServerDetails,
		Job:                job,
		LogLevel:           "ERROR", // DEBUG or ERROR
	}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(c interface{}) error {
	newConfig := c.(*Config)
	inst.config = newConfig
	return nil
}
