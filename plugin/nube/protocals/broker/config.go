package main

// Config is user plugin configuration
type Config struct {
	Port string `yaml:"port"`
}

// DefaultConfig implements plugin.Configurer
func (i *Instance) DefaultConfig() interface{} {
	return &Config{
		Port: "1882",
	}
}

// ValidateAndSetConfig implements plugin.Configurer
func (i *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	i.config = newConfig
	return nil
}
