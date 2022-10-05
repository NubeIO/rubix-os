package main

type Influx struct {
	Host        string  `yaml:"host"`
	Port        int     `yaml:"port"`
	Token       *string `yaml:"token"`
	Org         string  `yaml:"org"`
	Bucket      string  `yaml:"bucket"`
	Measurement string  `yaml:"measurement"`
}

type Job struct {
	Frequency string   `yaml:"frequency"`
	Networks  []string `yaml:"networks"`
}

type Config struct {
	Influx   []Influx `yaml:"influx"`
	Job      Job      `yaml:"job"`
	LogLevel string   `yaml:"log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	influx := Influx{
		Host:        "localhost",
		Port:        8086,
		Token:       nil,
		Org:         "nube-org",
		Bucket:      "nube-bucket",
		Measurement: "points",
	}
	job := Job{
		Frequency: "1m",
	}

	return &Config{
		Influx:   []Influx{influx},
		Job:      job,
		LogLevel: "ERROR", // DEBUG or ERROR
	}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	inst.config = newConfig
	return nil
}
