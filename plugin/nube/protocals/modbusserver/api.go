package main

import (
	"github.com/gin-gonic/gin"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
}
