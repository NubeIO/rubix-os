package main

// Config is user plugin configuration
type Config struct {
}

// DefaultConfig implements plugin.Configurer
func (inst *Instance) DefaultConfig() interface{} {
	return &Config{}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

// ValidateAndSetConfig implements plugin.Configurer
func (inst *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	inst.config = newConfig
	return nil
}
