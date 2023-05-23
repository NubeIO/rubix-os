package main

import "github.com/NubeIO/rubix-os/utils/float"

type Config struct {
	Job      Job    `yaml:"job"`
	LogLevel string `yaml:"log_level"`
}

type Job struct {
	Frequency                string            `yaml:"frequency"`
	FFHost                   string            `yaml:"ff_host"`
	FFPort                   float64           `yaml:"ff_port"`
	AlertDelayMins           *float64          `yaml:"alert_delay_mins"`
	SiteNamesInclude         []string          `yaml:"site_names_include"`
	SiteNamesExclude         []string          `yaml:"site_names_exclude"`
	RubixNetworkNamesInclude []string          `yaml:"rubix_network_names_include"`
	RubixNetworkNamesExclude []string          `yaml:"rubix_network_names_exclude"`
	RubixDeviceNamesInclude  []string          `yaml:"rubix_device_names_include"`
	RubixDeviceNamesExclude  []string          `yaml:"rubix_device_names_exclude"`
	RubixPointNamesInclude   []string          `yaml:"rubix_point_names_include"`
	RubixPointNamesExclude   []string          `yaml:"rubix_point_names_exclude"`
	TagsInclude              []string          `yaml:"tags_include"`
	TagsExclude              []string          `yaml:"tags_exclude"`
	MetaTagsInclude          map[string]string `yaml:"meta_tags_include"`
	MetaTagsExclude          map[string]string `yaml:"meta_tags_exclude"`
}

func (inst *Instance) DefaultConfig() interface{} {
	job := Job{
		Frequency:      "30m",
		FFHost:         "192.168.1.10",
		FFPort:         1616,
		AlertDelayMins: float.New(60),
	}

	return &Config{
		Job:      job,
		LogLevel: "DEBUG", // DEBUG or ERROR
	}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(c interface{}) error {
	newConfig := c.(*Config)
	inst.config = newConfig
	return nil
}
