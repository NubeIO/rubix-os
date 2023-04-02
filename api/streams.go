package api

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type StreamDatabase interface {
	GetStreams(args Args) ([]*model.Stream, error)
	GetStream(uuid string, args Args) (*model.Stream, error)
	CreateStream(body *model.Stream) (*model.Stream, error)
	UpdateStream(uuid string, body *model.Stream, checkAm bool) (*model.Stream, error)
	DeleteStream(uuid string) (bool, error)
	SyncStreamProducers(uuid string, args Args) ([]*interfaces.SyncModel, error)
}

type StreamAPI struct {
	DB StreamDatabase
}

func (j *StreamAPI) GetStreams(ctx *gin.Context) {
	args := buildStreamArgs(ctx)
	q, err := j.DB.GetStreams(args)
	ResponseHandler(q, err, ctx)
}

func (j *StreamAPI) GetStream(ctx *gin.Context) {
	args := buildStreamArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.GetStream(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *StreamAPI) CreateStream(ctx *gin.Context) {
	body, _ := getBODYStream(ctx)
	q, err := j.DB.CreateStream(body)
	ResponseHandler(q, err, ctx)
}

func (j *StreamAPI) UpdateStream(ctx *gin.Context) {
	body, _ := getBODYStream(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateStream(uuid, body, true)
	ResponseHandler(q, err, ctx)
}

func (j *StreamAPI) DeleteStream(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteStream(uuid)
	ResponseHandler(q, err, ctx)
}

func (j *StreamAPI) SyncStreamProducers(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildStreamArgs(ctx)
	q, err := j.DB.SyncStreamProducers(uuid, args)
	ResponseHandler(q, err, ctx)
}
