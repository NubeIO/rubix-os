package main

type Config struct {
	LogLevel string `yaml:"log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
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
