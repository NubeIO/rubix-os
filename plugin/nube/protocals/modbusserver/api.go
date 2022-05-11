package main

import (
	"github.com/gin-gonic/gin"
)

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
)

const (
	networks = "/networks"
	devices  = "/devices"
	points   = "/points"
)

var err error

// RegisterWebhook implements plugin.Webhooker
func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath

}
