package main

// Config is user plugin configuration
type Config struct {
	EnablePolling bool `yaml:"enable"`
}

// DefaultConfig implements plugin.Configurer
func (i *Instance) DefaultConfig() interface{} {
	return &Config{
		EnablePolling: false,
	}
}

// ValidateAndSetConfig implements plugin.Configurer
func (i *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	i.config = newConfig
	return nil
}
