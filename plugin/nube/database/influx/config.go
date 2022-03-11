package main

type Config struct {
	Influx Influx `yaml:"influx"`
}

type Influx struct {
	Host  string  `yaml:"host"`
	Port  int     `yaml:"port"`
	Token *string `yaml:"token"`
}

func (i *Instance) DefaultConfig() interface{} {
	influx := Influx{
		Host:  "localhost",
		Port:  8086,
		Token: nil,
	}
	return &Config{
		Influx: influx,
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
