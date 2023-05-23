package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type StreamCloneDatabase interface {
	GetStreamClones(args Args) ([]*model.StreamClone, error)
	GetStreamClone(uuid string, args Args) (*model.StreamClone, error)
	DeleteStreamClone(uuid string) (bool, error)
	DeleteOneStreamCloneByArgs(args Args) (bool, error)
	SyncStreamCloneConsumers(uuid string, args Args) ([]*interfaces.SyncModel, error)
}

type StreamCloneAPI struct {
	DB StreamCloneDatabase
}

func (j *StreamCloneAPI) GetStreamClones(ctx *gin.Context) {
	args := buildStreamCloneArgs(ctx)
	q, err := j.DB.GetStreamClones(args)
	ResponseHandler(q, err, ctx)
}

func (j *StreamCloneAPI) GetStreamClone(ctx *gin.Context) {
	args := buildStreamCloneArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.GetStreamClone(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *StreamCloneAPI) DeleteStreamClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteStreamClone(uuid)
	ResponseHandler(q, err, ctx)
}

func (j *StreamCloneAPI) DeleteOneStreamCloneByArgs(ctx *gin.Context) {
	args := buildStreamCloneArgs(ctx)
	q, err := j.DB.DeleteOneStreamCloneByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (j *StreamCloneAPI) SyncStreamCloneConsumers(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildStreamCloneArgs(ctx)
	q, err := j.DB.SyncStreamCloneConsumers(uuid, args)
	ResponseHandler(q, err, ctx)
}
