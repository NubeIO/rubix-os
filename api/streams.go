package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type StreamDatabase interface {
	GetStreams(args Args) ([]*model.Stream, error)
	GetStream(uuid string, args Args) (*model.Stream, error)
	CreateStream(body *model.Stream) (*model.Stream, error)
	UpdateStream(uuid string, body *model.Stream) (*model.Stream, error)
	DeleteStream(uuid string) (bool, error)
	DropStreams() (bool, error)
}

type StreamAPI struct {
	DB StreamDatabase
}

func (j *StreamAPI) GetStreams(ctx *gin.Context) {
	args := buildStreamArgs(ctx)
	q, err := j.DB.GetStreams(args)
	responseHandler(q, err, ctx)
}

func (j *StreamAPI) GetStream(ctx *gin.Context) {
	args := buildStreamArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.GetStream(uuid, args)
	responseHandler(q, err, ctx)
}

func (j *StreamAPI) CreateStream(ctx *gin.Context) {
	body, _ := getBODYStream(ctx)
	q, err := j.DB.CreateStream(body)
	responseHandler(q, err, ctx)
}

func (j *StreamAPI) UpdateStream(ctx *gin.Context) {
	body, _ := getBODYStream(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateStream(uuid, body)
	responseHandler(q, err, ctx)
}

func (j *StreamAPI) DeleteStream(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteStream(uuid)
	responseHandler(q, err, ctx)
}

func (j *StreamAPI) DropStreams(ctx *gin.Context) {
	q, err := j.DB.DropStreams()
	responseHandler(q, err, ctx)
}
