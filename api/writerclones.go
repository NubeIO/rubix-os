package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The WriterCloneDatabase interface for encapsulating database access.
type WriterCloneDatabase interface {
	GetWriterClone(uuid string) (*model.WriterClone, error)
	GetWriterClones() ([]*model.WriterClone, error)
	CreateWriterClone(body *model.WriterClone) (*model.WriterClone, error)
	UpdateWriterClone(uuid string, body *model.WriterClone, updateProducer bool) (*model.WriterClone, error)
	DeleteWriterClone(uuid string) (bool, error)
	DropWriterClone() (bool, error)
}
type WriterCloneAPI struct {
	DB WriterCloneDatabase
}

func (j *WriterCloneAPI) GetWriterClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetWriterClone(uuid)
	reposeHandler(q, err, ctx)
}

func (j *WriterCloneAPI) GetWriterClones(ctx *gin.Context) {
	q, err := j.DB.GetWriterClones()
	reposeHandler(q, err, ctx)

}

func (j *WriterCloneAPI) CreateWriterClone(ctx *gin.Context) {
	body, _ := getBODYWriterClone(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateWriterClone(body)
	reposeHandler(q, err, ctx)
}

func (j *WriterCloneAPI) UpdateWriterClone(ctx *gin.Context) {
	body, _ := getBODYWriterClone(ctx)
	uuid := resolveID(ctx)
	_, _, _, updateProducer := withConsumerArgs(ctx)
	q, err := j.DB.UpdateWriterClone(uuid, body, updateProducer)
	reposeHandler(q, err, ctx)
}

func (j *WriterCloneAPI) DeleteWriterClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteWriterClone(uuid)
	reposeHandler(q, err, ctx)
}

func (j *WriterCloneAPI) DropWriterClone(ctx *gin.Context) {
	q, err := j.DB.DropWriterClone()
	reposeHandler(q, err, ctx)
}
