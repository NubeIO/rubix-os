package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type WriterDatabase interface {
	GetWriters(args Args) ([]*model.Writer, error)
	GetWriter(uuid string) (*model.Writer, error)
	GetWriterByName(flowNetworkCloneName string, streamCloneName string, consumerName string,
		writerThingName string) (*model.Writer, error)
	CreateWriter(body *model.Writer) (*model.Writer, error)
	UpdateWriter(uuid string, body *model.Writer, checkAm bool) (*model.Writer, error)
	DeleteWriter(uuid string) (bool, error)
	WriterAction(uuid string, body *model.WriterBody) *model.WriterActionOutput
	WriterBulkAction(body []*model.WriterBulkBody) []*model.WriterActionOutput
}

type WriterAPI struct {
	DB WriterDatabase
}

func (j *WriterAPI) GetWriters(ctx *gin.Context) {
	args := buildWriterArgs(ctx)
	q, err := j.DB.GetWriters(args)
	ResponseHandler(q, err, ctx)
}

func (j *WriterAPI) GetWriter(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetWriter(uuid)
	ResponseHandler(q, err, ctx)
}

func (j *WriterAPI) GetWriterByName(ctx *gin.Context) {
	flowNetworkCloneName := resolveFlowNetworkCloneName(ctx)
	streamCloneName := resolveStreamCloneName(ctx)
	consumerName := resolveConsumerName(ctx)
	writerThingName := resolveWriterThingName(ctx)
	q, err := j.DB.GetWriterByName(flowNetworkCloneName, streamCloneName, consumerName, writerThingName)
	ResponseHandler(q, err, ctx)
}

func (j *WriterAPI) CreateWriter(ctx *gin.Context) {
	body, _ := getBODYWriter(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateWriter(body)
	ResponseHandler(q, err, ctx)
}

func (j *WriterAPI) UpdateWriter(ctx *gin.Context) {
	body, _ := getBODYWriter(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateWriter(uuid, body, true)
	ResponseHandler(q, err, ctx)
}

func (j *WriterAPI) DeleteWriter(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteWriter(uuid)
	ResponseHandler(q, err, ctx)
}

func (j *WriterAPI) WriterAction(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBODYWriterBody(ctx)
	q := j.DB.WriterAction(uuid, body)
	ResponseHandler(q, nil, ctx)
}

func (j *WriterAPI) WriterBulkAction(ctx *gin.Context) {
	body, _ := getBODYWriterBulk(ctx)
	q := j.DB.WriterBulkAction(body)
	ResponseHandler(q, nil, ctx)
}

type WriterWizard struct {
	ConsumerFlowUUID   string `json:"consumer_side_flow_uuid"`
	ConsumerStreamUUID string `json:"consumer_side_stream_uuid"`
	ProducerUUID       string `json:"remote_producer_uuid"`
}
