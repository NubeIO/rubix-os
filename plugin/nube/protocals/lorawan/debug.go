package main

import (
	"github.com/NubeIO/flow-framework/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) lorawanDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "Lorawan: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) lorawanErrorMsg(args ...interface{}) {
	prefix := "Lorawan: "
	log.Error(prefix, args)
}
