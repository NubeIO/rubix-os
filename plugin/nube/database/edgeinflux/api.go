package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/gin-gonic/gin"
)

const (
	help = "/help"
)

func (inst *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	inst.basePath = basePath
	inst.edgeinfluxDebugMsg(fmt.Sprintf("RegisterWebhook(): %+v\n", inst))
	mux.PATCH(plugin.PointsWriteURL, func(ctx *gin.Context) {
		uuid := plugin.ResolveID(ctx)
		err := inst.SendPointWriteHistory(uuid)
		api.ResponseHandler(nil, err, ctx)
	})
}
