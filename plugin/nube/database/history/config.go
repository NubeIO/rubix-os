package main

type Job struct {
	Frequency string `yaml:"frequency"`
}

type Config struct {
	Job Job `yaml:"job"`
}

func (inst *Instance) DefaultConfig() interface{} {
	job := Job{
		Frequency: "1m",
	}
	return &Config{
		Job: job,
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
