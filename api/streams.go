package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The StreamDatabase interface for encapsulating database access.
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
	args := streamArgs(ctx)
	q, err := j.DB.GetStreams(args)
	reposeHandler(q, err, ctx)
}

func (j *StreamAPI) GetStream(ctx *gin.Context) {
	args := streamArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.GetStream(uuid, args)
	reposeHandler(q, err, ctx)
}

func (j *StreamAPI) CreateStream(ctx *gin.Context) {
	body, _ := getBODYStream(ctx)
	q, err := j.DB.CreateStream(body)
	reposeHandler(q, err, ctx)
}

func (j *StreamAPI) UpdateStream(ctx *gin.Context) {
	body, _ := getBODYStream(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateStream(uuid, body)
	reposeHandler(q, err, ctx)
}

func (j *StreamAPI) DeleteStream(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteStream(uuid)
	reposeHandler(q, err, ctx)
}

func (j *StreamAPI) DropStreams(ctx *gin.Context) {
	q, err := j.DB.DropStreams()
	reposeHandler(q, err, ctx)
}
