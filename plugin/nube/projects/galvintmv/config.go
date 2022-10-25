package main

type Config struct {
	Job      Job    `yaml:"job"`
	LogLevel string `yaml:"log_level"`
}

type Job struct {
	EnableConfigSteps           bool    `yaml:"enable_config_steps"`
	Frequency                   string  `yaml:"frequency"`
	ChirpstackHost              string  `yaml:"chirpstack_host"`
	ChirpstackPort              float64 `yaml:"chirpstack_port"`
	ChirpstackApplicationNumber int     `yaml:"chirpstack_application_number"`
	ChirpstackNetworkKey        string  `yaml:"chirpstack_network_key"`
	ChirpstackUsername          string  `yaml:"chirpstack_username"`
	ChirpstackPassword          string  `yaml:"chirpstack_password"`
	DeviceJSONFilePath          string  `yaml:"device_json_file_path"`
}

func (inst *Instance) DefaultConfig() interface{} {
	job := Job{
		EnableConfigSteps:           false,
		Frequency:                   "30m",
		ChirpstackHost:              "0.0.0.0",
		ChirpstackPort:              8080,
		ChirpstackApplicationNumber: 1,
		DeviceJSONFilePath:          "/home/pi/test.json",
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
