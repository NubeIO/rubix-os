package main

import (
	"github.com/NubeIO/rubix-os/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) edgeazureDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "edgeAzure: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) edgeazureErrorMsg(args ...interface{}) {
	prefix := "edgeAzure: "
	log.Error(prefix, args)
}
