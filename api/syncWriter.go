package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type SyncWriterDatabase interface {
	SyncWriter(body *model.SyncWriter) (*model.WriterClone, error)
	SyncCOV(writerUUID string, body *model.SyncCOV) error
	SyncWriterWriteAction(sourceUUID string, body *model.SyncWriterAction) error
	SyncWriterReadAction(sourceUUID string) error
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
	writerUUID := resolveWriterUUID(ctx)
	body, _ := getBodySyncCOV(ctx)
	err := a.DB.SyncCOV(writerUUID, body)
	responseHandler(nil, err, ctx)
}

func (a *SyncWriterAPI) SyncWriterWriteAction(ctx *gin.Context) {
	sourceUUID := resolveSourceUUID(ctx)
	body, _ := getBodySyncWriterAction(ctx)
	err := a.DB.SyncWriterWriteAction(sourceUUID, body)
	responseHandler(nil, err, ctx)
}

func (a *SyncWriterAPI) SyncWriterReadAction(ctx *gin.Context) {
	sourceUUID := resolveSourceUUID(ctx)
	err := a.DB.SyncWriterReadAction(sourceUUID)
	responseHandler(nil, err, ctx)
}
