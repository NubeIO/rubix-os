package main

type Config struct {
	EnablePolling bool `yaml:"enable_polling"`
}

func (i *Instance) DefaultConfig() interface{} {
	return &Config{
		EnablePolling: true,
	}
}

func (i *Instance) GetConfig() interface{} {
	return i.config
}

func (i *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	i.config = newConfig
	return nil
}
