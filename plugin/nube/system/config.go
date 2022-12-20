package main

type Config struct {
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	inst.config = newConfig
	return nil
}
