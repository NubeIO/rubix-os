package main

type Azure struct {
	HostName        string `json:"azure_host_name"`
	DeviceId        string `json:"azure_device_id"`
	SharedAccessKey string `json:"azure_shared_access_key"`
}

type Job struct {
	Frequency            string   `yaml:"frequency"`
	Networks             []string `yaml:"networks"`
	RequireHistoryEnable bool     `yaml:"require_history_enable"`
}

type Config struct {
	Azure    Azure  `yaml:"azure"`
	Job      Job    `yaml:"job"`
	LogLevel string `yaml:"log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	azure := Azure{
		HostName:        "",
		DeviceId:        "",
		SharedAccessKey: "",
	}
	job := Job{
		Frequency:            "60m",
		Networks:             []string{"system"},
		RequireHistoryEnable: true,
	}

	return &Config{
		Azure:    azure,
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
