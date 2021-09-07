package main

// Config is user plugin configuration
type Config struct {
	MagicString string `yaml:"magic_string"`
}

// DefaultConfig implements plugin.Configurer
func (c *Instance) DefaultConfig() interface{} {
	return &Config{
		MagicString: "hello world",
	}
}

// ValidateAndSetConfig implements plugin.Configurer
func (c *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	c.config = newConfig
	return nil
}
