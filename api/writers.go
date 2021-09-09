package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The WriterDatabase interface for encapsulating database access.
type WriterDatabase interface {
	GetWriter(uuid string) (*model.Writer, error)
	GetWriters() ([]*model.Writer, error)
	CreateWriter(body *model.Writer) (*model.Writer, error)
	UpdateWriter(uuid string, body *model.Writer) (*model.Writer, error)
	DeleteWriter(uuid string) (bool, error)
	WriterAction(uuid string, body *model.WriterBody) (*model.ProducerHistory, error)
	WriterBulkAction(body []*model.WriterBulk) ([]*model.ProducerHistory, error)
}

type WriterAPI struct {
	DB WriterDatabase
}

func (j *WriterAPI) GetWriter(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetWriter(uuid)
	reposeHandler(q, err, ctx)
}

func (j *WriterAPI) GetWriters(ctx *gin.Context) {
	q, err := j.DB.GetWriters()
	reposeHandler(q, err, ctx)

}

func (j *WriterAPI) CreateWriter(ctx *gin.Context) {
	body, _ := getBODYWriter(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateWriter(body)
	reposeHandler(q, err, ctx)
}

func (j *WriterAPI) UpdateWriter(ctx *gin.Context) {
	body, _ := getBODYWriter(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateWriter(uuid, body)
	reposeHandler(q, err, ctx)
}

func (j *WriterAPI) DeleteWriter(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteWriter(uuid)
	reposeHandler(q, err, ctx)
}

//WriterAction get or update a producer value by using the writer uuid
func (j *WriterAPI) WriterAction(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBODYWriterBody(ctx)
	q, err := j.DB.WriterAction(uuid, body)
	reposeHandler(q, err, ctx)
}

//WriterBulkAction get or update a producer value by using the writer uuid
func (j *WriterAPI) WriterBulkAction(ctx *gin.Context) {
	body, _ := getBODYWriterBulk(ctx)
	q, err := j.DB.WriterBulkAction(body)
	reposeHandler(q, err, ctx)
}
