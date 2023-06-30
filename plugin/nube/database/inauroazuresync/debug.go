package main

import (
	"github.com/NubeIO/rubix-os/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) inauroazuresyncDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "inauroazuresync: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) inauroazuresyncErrorMsg(args ...interface{}) {
	prefix := "inauroazuresync: "
	log.Error(prefix, args)
}
