package main

// Config is user plugin configuration
type Config struct {
	Host            string `yaml:"host"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

// DefaultConfig implements plugin.Configurer
func (i *Instance) DefaultConfig() interface{} {
	return &Config{
		Host:            "127.0.0.1:9000",
		AccessKeyID:     "12345678",
		SecretAccessKey: "12345678",
	}
}

// ValidateAndSetConfig implements plugin.Configurer
func (i *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	i.config = newConfig
	return nil
}
