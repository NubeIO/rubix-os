package main

type Config struct {
	EnablePolling     bool   `yaml:"enable_polling"`
	LogLevel          string `yaml:"log_level"`
	PollQueueLogLevel string `yaml:"poll_queue_log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		EnablePolling:     true,
		LogLevel:          "ERROR", // DEBUG or ERROR
		PollQueueLogLevel: "ERROR", // DEBUG or ERROR
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
