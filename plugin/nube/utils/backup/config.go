package main

type Config struct {
	Host            string `yaml:"host"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	BucketName      string `yaml:"bucket_name"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		Host:            "127.0.0.1:9000",
		AccessKeyID:     "",
		SecretAccessKey: "",
		BucketName:      "",
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
