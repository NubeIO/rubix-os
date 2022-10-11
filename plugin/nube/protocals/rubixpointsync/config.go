package main

type Job struct {
	Frequency           string   `yaml:"frequency"`
	Networks            []string `yaml:"networks"`
	RequireNetworkMatch bool     `yaml:"require_network_match"`
	GenerateRubixPoints bool     `yaml:"generate_rubix_points"`
}

type Config struct {
	Job      Job    `yaml:"job"`
	LogLevel string `yaml:"log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	job := Job{
		Frequency:           "15m",
		Networks:            []string{"system"},
		RequireNetworkMatch: true,
		GenerateRubixPoints: false,
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
