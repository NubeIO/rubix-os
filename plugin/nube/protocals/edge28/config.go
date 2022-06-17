package main

import "github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/config"

func (inst *Instance) DefaultConfig() interface{} {
	return &config.Config{
		EnablePolling: true,
		LogLevel:      "ERROR", // DEBUG or ERROR
	}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(c interface{}) error {
	newConfig := c.(*config.Config)
	inst.config = newConfig
	return nil
}
