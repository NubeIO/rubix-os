package main

type Config struct {
	EnablePolling bool   `yaml:"enable_polling"`
	Ip            string `yaml:"ip"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		EnablePolling: true,
		Ip:            "0.0.0.0",
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
