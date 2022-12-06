package main

import (
	"github.com/NubeIO/flow-framework/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) flatlinealertsDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "Flatline Alerts: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) flatlinealertsErrorMsg(args ...interface{}) {
	prefix := "Flatline Alerts: "
	log.Error(prefix, args)
}
