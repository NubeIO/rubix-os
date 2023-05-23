package main

import (
	"github.com/NubeIO/rubix-os/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) statusmismatchalertsDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "StatusMismatch Alerts: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) statusmismatchalertsErrorMsg(args ...interface{}) {
	prefix := "StatusMismatch Alerts: "
	log.Error(prefix, args)
}
