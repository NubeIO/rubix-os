package main

import (
	"github.com/NubeIO/flow-framework/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) rubixpointsyncDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "pointserversync: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) rubixpointsyncErrorMsg(args ...interface{}) {
	prefix := "pointserversync: "
	log.Error(prefix, args)
}
