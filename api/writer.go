package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type WriterDatabase interface {
	GetWriter(uuid string) (*model.Writer, error)
	GetWriters(args Args) ([]*model.Writer, error)
	CreateWriter(body *model.Writer) (*model.Writer, error)
	UpdateWriter(uuid string, body *model.Writer) (*model.Writer, error)
	DeleteWriter(uuid string) (bool, error)
	DropWriters() (bool, error)
	WriterAction(uuid string, body *model.WriterBody) *model.WriterActionOutput
	WriterBulkAction(body []*model.WriterBulkBody) []*model.WriterActionOutput
}

type WriterAPI struct {
	DB WriterDatabase
}

func (j *WriterAPI) GetWriter(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetWriter(uuid)
	responseHandler(q, err, ctx)
}

func (j *WriterAPI) GetWriters(ctx *gin.Context) {
	args := buildWriterArgs(ctx)
	q, err := j.DB.GetWriters(args)
	responseHandler(q, err, ctx)
}

func (j *WriterAPI) CreateWriter(ctx *gin.Context) {
	body, _ := getBODYWriter(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		responseHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateWriter(body)
	responseHandler(q, err, ctx)
}

func (j *WriterAPI) UpdateWriter(ctx *gin.Context) {
	body, _ := getBODYWriter(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateWriter(uuid, body)
	responseHandler(q, err, ctx)
}

func (j *WriterAPI) DeleteWriter(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteWriter(uuid)
	responseHandler(q, err, ctx)
}

func (j *WriterAPI) DropWriters(ctx *gin.Context) {
	q, err := j.DB.DropWriters()
	responseHandler(q, err, ctx)
}

func (j *WriterAPI) WriterAction(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBODYWriterBody(ctx)
	q := j.DB.WriterAction(uuid, body)
	responseHandler(q, nil, ctx)
}

func (j *WriterAPI) WriterBulkAction(ctx *gin.Context) {
	body, _ := getBODYWriterBulk(ctx)
	q := j.DB.WriterBulkAction(body)
	responseHandler(q, nil, ctx)
}

type WriterWizard struct {
	ConsumerFlowUUID   string `json:"consumer_side_flow_uuid"`
	ConsumerStreamUUID string `json:"consumer_side_stream_uuid"`
	ProducerUUID       string `json:"remote_producer_uuid"`
}
