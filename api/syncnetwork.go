package api

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type SyncNetworkDatabase interface {
	SyncNetwork(fn *interfaces.SyncNetwork) (*model.Network, error)
}

type SyncNetworkAPI struct {
	DB SyncNetworkDatabase
}

func (a *SyncNetworkAPI) SyncNetwork(ctx *gin.Context) {
	body, _ := getBodySyncNetwork(ctx)
	q, err := a.DB.SyncNetwork(body)
	ResponseHandler(q, err, ctx)
}
