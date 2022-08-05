package main

import (
	"github.com/gin-gonic/gin"
)

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"
	networks      = "/networks"
	devices       = "/devices"
	points        = "/points"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
}
