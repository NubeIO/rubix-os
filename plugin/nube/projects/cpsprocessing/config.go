package main

type Config struct {
	Postgres Postgres `yaml:"postgres"`
	Job      Job      `yaml:"job"`
	LogLevel string   `yaml:"log_level"`
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
	Frequency string `yaml:"frequency"`
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
	job := Job{
		Frequency: "1m",
	}

	return &Config{
		Postgres: postgres,
		Job:      job,
		LogLevel: "ERROR", // DEBUG or ERROR
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
