package main

// Config is user plugin configuration
type Config struct {
	MagicString string `yaml:"magic_string"`
}

// DefaultConfig implements plugin.Configurer
func (i *Instance) DefaultConfig() interface{} {
	return &Config{
		MagicString: "N/A",
	}
}

func (i *Instance) GetConfig() interface{} {
	return i.config
}

// ValidateAndSetConfig implements plugin.Configurer
func (i *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	i.config = newConfig
	return nil
}
