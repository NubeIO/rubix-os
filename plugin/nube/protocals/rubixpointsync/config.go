package main

type Job struct {
	Host                string   `yaml:"host"`
	Port                float64  `yaml:"port"`
	Frequency           string   `yaml:"frequency"`
	Networks            []string `yaml:"networks"`
	RequireNetworkMatch bool     `yaml:"require_network_match"`
	GenerateRubixPoints bool     `yaml:"generate_rubix_points"`
	UpdateOnCOV         bool     `yaml:"update_on_cov"`
}

type Config struct {
	Job      Job    `yaml:"job"`
	LogLevel string `yaml:"log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	job := Job{
		Host:                "0.0.0.0",
		Port:                1515,
		Frequency:           "15m",
		Networks:            []string{"system"},
		RequireNetworkMatch: true,
		GenerateRubixPoints: false,
		UpdateOnCOV:         true,
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
