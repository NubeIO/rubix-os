package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type SyncStreamDatabase interface {
	SyncStream(fn *model.StreamSync) (*model.StreamClone, error)
}

type SyncStreamAPI struct {
	DB SyncStreamDatabase
}

func (a *SyncStreamAPI) SyncStream(ctx *gin.Context) {
	body, _ := getBodyStreamSync(ctx)
	q, err := a.DB.SyncStream(body)
	reposeHandler(q, err, ctx)
}
