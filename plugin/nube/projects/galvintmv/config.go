package main

type Config struct {
	Job      Job    `yaml:"job"`
	LogLevel string `yaml:"log_level"`
}

type Job struct {
	Frequency      string   `yaml:"frequency"`
	Networks       []string `yaml:"networks"`
	ChirpstackHost string   `yaml:"host"`
	ChirpstackPort float64  `yaml:"port"`
}

func (inst *Instance) DefaultConfig() interface{} {
	job := Job{
		Frequency:      "60m",
		ChirpstackHost: "0.0.0.0",
		ChirpstackPort: 8080,
	}

	return &Config{
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
