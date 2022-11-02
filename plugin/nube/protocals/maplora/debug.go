package main

import (
	"github.com/NubeIO/flow-framework/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) maploraDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "maplora: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) maploraErrorMsg(args ...interface{}) {
	prefix := "maplora: "
	log.Error(prefix, args)
}
