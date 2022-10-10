package main

type Config struct {
	Writers                  []string
	ValueSyncIntervalSeconds int
	LinkSyncIntervalSeconds  int
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		Writers:                  []string{"modbus", "bacnetmaster"},
		ValueSyncIntervalSeconds: 20,
		LinkSyncIntervalSeconds:  120,
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
