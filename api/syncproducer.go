package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type SyncProducerDatabase interface {
	SyncProducer(fn *interfaces.SyncProducer) ([]*model.Consumer, error)
}

type SyncProducerAPI struct {
	DB SyncProducerDatabase
}

func (a *SyncProducerAPI) SyncProducer(ctx *gin.Context) {
	body, _ := getBodySyncProducer(ctx)
	q, err := a.DB.SyncProducer(body)
	ResponseHandler(q, err, ctx)
}
