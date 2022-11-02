package main

import (
	"github.com/NubeIO/flow-framework/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) mapmodbusDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "mapmodbus: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) mapmodbusErrorMsg(args ...interface{}) {
	prefix := "mapmodbus: "
	log.Error(prefix, args)
}
