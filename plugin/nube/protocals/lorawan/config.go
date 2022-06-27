package main

type Config struct {
	CSAddress            string `default:"localhost"`
	CSPort               int    `default:"8080"`
	CSUsername           string `default:"admin"`
	CSPassword           string `default:"admin"`
	DeviceLimit          int    `default:"200"`
	SyncPeriodMins       int    `default:"1"`
	ReconnectTimeoutSecs int    `default:"10"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		CSAddress:            "0.0.0.0",
		CSPort:               8080,
		CSUsername:           "admin",
		CSPassword:           "N00BWAN",
		DeviceLimit:          200,
		SyncPeriodMins:       1,
		ReconnectTimeoutSecs: 10,
	}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := inst.DefaultConfig().(*Config)
	inst.config = newConfig
	return nil
}
