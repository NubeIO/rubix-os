package main

import (
	"github.com/NubeIO/rubix-os/utils/nstring"
	log "github.com/sirupsen/logrus"
)

func (inst *Instance) edgeinfluxDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(inst.config.LogLevel, "DEBUG") {
		prefix := "edgeInflux: "
		log.Info(prefix, args)
	}
}

func (inst *Instance) edgeinfluxErrorMsg(args ...interface{}) {
	prefix := "edgeInflux: "
	log.Error(prefix, args)
}
