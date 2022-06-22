package main

type Config struct {
	CSAddress         string `default:"localhost"`
	CSPort            int    `default:"8080"`
	CSUsername        string `default:"admin"`
	CSPassword        string `default:"admin"`
	DeviceLimit       int    `default:"200"`
	SyncPeriodMinutes int    `default:"1"`
}

func (inst *Instance) setTemporaryConfigDefaults() {
	inst.config.CSAddress = "192.168.1.114"
	inst.config.CSPort = 8080
	inst.config.CSUsername = "admin"
	inst.config.CSPassword = "admin"
	inst.config.DeviceLimit = 200
	inst.config.SyncPeriodMinutes = 1
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	inst.config = newConfig
	inst.setTemporaryConfigDefaults()
	return nil
}
