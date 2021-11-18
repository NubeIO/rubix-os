package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type SyncStreamDatabase interface {
	SyncStream(fn *model.SyncStream) (*model.StreamClone, error)
}

type SyncStreamAPI struct {
	DB SyncStreamDatabase
}

func (a *SyncStreamAPI) SyncStream(ctx *gin.Context) {
	body, _ := getBodySyncStream(ctx)
	q, err := a.DB.SyncStream(body)
	responseHandler(q, err, ctx)
}
