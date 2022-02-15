package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type SyncWriterDatabase interface {
	SyncWriter(body *model.SyncWriter) (*model.WriterClone, error)
	SyncCOV(body *model.SyncCOV) error
	SyncWriterAction(body *model.SyncWriterAction) error
}

type SyncWriterAPI struct {
	DB SyncWriterDatabase
}

func (a *SyncWriterAPI) SyncWriter(ctx *gin.Context) {
	body, _ := getBodySyncWriter(ctx)
	q, err := a.DB.SyncWriter(body)
	responseHandler(q, err, ctx)
}

func (a *SyncWriterAPI) SyncCOV(ctx *gin.Context) {
	body, _ := getBodySyncCOV(ctx)
	err := a.DB.SyncCOV(body)
	responseHandler(nil, err, ctx)
}

func (a *SyncWriterAPI) SyncWriterAction(ctx *gin.Context) {
	body, _ := getBodySyncWriterAction(ctx)
	err := a.DB.SyncWriterAction(body)
	responseHandler(nil, err, ctx)
}
