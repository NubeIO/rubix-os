package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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
	ResponseHandler(q, err, ctx)
}
