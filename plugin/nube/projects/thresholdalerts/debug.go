package main

import (
	"github.com/NubeIO/flow-framework/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) thresholdalertsDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "Threshold Alerts: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) thresholdalertsErrorMsg(args ...interface{}) {
	prefix := "Threshold Alerts: "
	log.Error(prefix, args)
}
