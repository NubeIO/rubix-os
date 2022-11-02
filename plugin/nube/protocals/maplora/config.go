package main

type Job struct {
	Host               string  `yaml:"host"`
	Port               float64 `yaml:"port"`
	UseExistingNetwork bool    `yaml:"use_existing_network"`
}

type Config struct {
	Job      Job    `yaml:"job"`
	LogLevel string `yaml:"log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	job := Job{
		Host:               "0.0.0.0",
		Port:               1616,
		UseExistingNetwork: false,
	}

	return &Config{
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
