package main

import (
	"github.com/NubeIO/flow-framework/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) systemDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "System: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) systemErrorMsg(args ...interface{}) {
	prefix := "System: "
	log.Error(prefix, args)
}
