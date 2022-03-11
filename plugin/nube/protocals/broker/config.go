package main

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	HttpPort string `yaml:"http_port"`
}

func (i *Instance) DefaultConfig() interface{} {
	return &Config{
		Host:     "0.0.0.0",
		Port:     "1883",
		HttpPort: "8099",
	}
}

func (i *Instance) GetConfig() interface{} {
	return i.config
}

func (i *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	i.config = newConfig
	return nil
}
