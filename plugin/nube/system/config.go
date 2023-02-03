package main

type Config struct {
	Schedule Schedule `yaml:"schedule"`
	LogLevel string   `yaml:"log_level"`
	// ScheduleLogLevel string `yaml:"schedule_log_level"`
}

type Schedule struct {
	Frequency string `yaml:"frequency"`
}

func (inst *Instance) DefaultConfig() interface{} {
	schedule := Schedule{
		Frequency: "60s",
	}

	return &Config{
		Schedule: schedule,
		LogLevel: "ERROR", // DEBUG or ERROR
		// ScheduleLogLevel: "ERROR", // DEBUG or ERROR
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
