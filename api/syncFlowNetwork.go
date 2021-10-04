package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type SyncFlowNetworkDatabase interface {
	SyncFlowNetwork(fn *model.FlowNetwork) (*model.FlowNetworkClone, error)
}

type SyncFlowNetworkAPI struct {
	DB SyncFlowNetworkDatabase
}

func (a *SyncFlowNetworkAPI) SyncFlowNetwork(ctx *gin.Context) {
	body, _ := getBODYFlowNetwork(ctx)
	q, err := a.DB.SyncFlowNetwork(body)
	reposeHandler(q, err, ctx)
}
