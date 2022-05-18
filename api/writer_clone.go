package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The WriterCloneDatabase interface for encapsulating database access.
type WriterCloneDatabase interface {
	GetWriterClones(args Args) ([]*model.WriterClone, error)
	GetWriterClone(uuid string) (*model.WriterClone, error)
	GetOneWriterCloneByArgs(args Args) (*model.WriterClone, error)
	CreateWriterClone(body *model.WriterClone) (*model.WriterClone, error)
	DeleteWriterClone(uuid string) (bool, error)
	DeleteOneWriterCloneByArgs(args Args) (bool, error)
}
type WriterCloneAPI struct {
	DB WriterCloneDatabase
}

func (j *WriterCloneAPI) GetWriterClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetWriterClone(uuid)
	responseHandler(q, err, ctx)
}

func (j *WriterCloneAPI) GetWriterClones(ctx *gin.Context) {
	args := buildWriterCloneArgs(ctx)
	q, err := j.DB.GetWriterClones(args)
	responseHandler(q, err, ctx)
}

func (j *WriterCloneAPI) CreateWriterClone(ctx *gin.Context) {
	body, _ := getBODYWriterClone(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		responseHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateWriterClone(body)
	responseHandler(q, err, ctx)
}

func (j *WriterCloneAPI) DeleteWriterClone(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteWriterClone(uuid)
	responseHandler(q, err, ctx)
}

func (j *WriterCloneAPI) DeleteOneWriterCloneByArgs(ctx *gin.Context) {
	args := buildWriterCloneArgs(ctx)
	q, err := j.DB.DeleteOneWriterCloneByArgs(args)
	responseHandler(q, err, ctx)
}
