package main

type Config struct {
	EnablePolling   bool   `yaml:"enable_polling"`
	Ip              string `yaml:"ip"`
	PollingTimeInMs int    `yaml:"polling_time_in_ms"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		EnablePolling:   true,
		Ip:              "0.0.0.0",
		PollingTimeInMs: 500,
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
