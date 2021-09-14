package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

/*
Stream
*/

// The StreamDatabase interface for encapsulating database access.
type StreamDatabase interface {
	GetStream(uuid string, withChildren bool) (*model.Stream, error)
	GetStreams(withChildren bool) ([]*model.Stream, error)
	CreateStream(body *model.Stream, AddToParent string) (*model.Stream, error)
	UpdateStream(uuid string, body *model.Stream) (*model.Stream, error)
	DeleteStream(uuid string) (bool, error)
	DropStreams() (bool, error)
}

type StreamAPI struct {
	DB StreamDatabase
}

func (j *StreamAPI) GetStreams(ctx *gin.Context) {
	withChildren, _, _ := withChildrenArgs(ctx)
	q, err := j.DB.GetStreams(withChildren)
	reposeHandler(q, err, ctx)
}

func (j *StreamAPI) GetStream(ctx *gin.Context) {
	withChildren, _, _ := withChildrenArgs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.GetStream(uuid, withChildren)
	reposeHandler(q, err, ctx)
}

func (j *StreamAPI) CreateStream(ctx *gin.Context) {
	body, _ := getBODYStream(ctx)
	AddToParent := parentArgs(ctx) //flowUUID
	q, err := j.DB.CreateStream(body, AddToParent)
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
