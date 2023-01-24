package main

type Config struct {
	CSAddress            string  `yaml:"csaddress"`
	CSPort               int     `yaml:"csport"`
	CSToken              string  `yaml:"cstoken"`
	CSUsername           string  `yaml:"csusername"`
	CSPassword           string  `yaml:"cspassword"`
	DeviceLimit          int     `yaml:"devicelimit"`
	SyncPeriodMins       float32 `yaml:"syncperiodminutes"`
	ReconnectTimeoutSecs int     `yaml:"reconnecttimeoutseconds"`
	LogLevel             string  `yaml:"log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		CSAddress:            "0.0.0.0",
		CSPort:               8080,
		CSToken:              "",
		CSUsername:           "",
		CSPassword:           "",
		DeviceLimit:          200,
		SyncPeriodMins:       1,
		ReconnectTimeoutSecs: 10,
		LogLevel:             "ERROR", // DEBUG or ERROR
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
