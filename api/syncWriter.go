package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type SyncWriterDatabase interface {
	SyncWriter(body *model.SyncWriter) (*model.WriterClone, error)
	SyncWriterCOV(body *model.SyncWriterCOV) error
}

type SyncWriterAPI struct {
	DB SyncWriterDatabase
}

func (a *SyncWriterAPI) SyncWriter(ctx *gin.Context) {
	body, _ := getBodySyncWriter(ctx)
	q, err := a.DB.SyncWriter(body)
	responseHandler(q, err, ctx)
}

func (a *SyncWriterAPI) SyncWriterCOV(ctx *gin.Context) {
	body, _ := getBodySyncWriterCOV(ctx)
	err := a.DB.SyncWriterCOV(body)
	responseHandler(nil, err, ctx)
}
