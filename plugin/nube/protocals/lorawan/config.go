package main

type Config struct {
	CSAddress            string  `yaml:"csaddress"`
	CSPort               int     `yaml:"csport"`
	CSToken              string  `yaml:"cstoken"`
	DeviceLimit          int     `yaml:"devicelimit"`
	SyncPeriodMins       float32 `yaml:"syncperiodminutes"`
	ReconnectTimeoutSecs int     `yaml:"reconnecttimeoutseconds"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		CSAddress:            "0.0.0.0",
		CSPort:               8080,
		CSToken:              "",
		DeviceLimit:          200,
		SyncPeriodMins:       1,
		ReconnectTimeoutSecs: 10,
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
