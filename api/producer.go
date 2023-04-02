package api

import (
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

/*
add producer
- user needs to pass a valid producer uuid (for example a point uuid) and type model.ProducerType (Point, Job, Alarm) and  model.ProducerApplication (Plugin, Remote, Local)

example for workflow for a point (Point 1 Has a Consumer to Point 2):
Point 1
-- consumer table -> point 2 uuid
-- producer table -> nil

Point 2
-- consumer table -> nil
-- producer table -> point 1 uuid

remote point will subscribe to cov events
- the local db will store a copy of the producer to know where to publish the data to
- the remote device will store a copy of its consumer in the consumers table, these will be the details of the remote producer

remote producer
- required: rubix-uuid
- optional: point uuid (required network_name and device_name and point_name)
*/

type ProducerDatabase interface {
	GetProducers(args Args) ([]*model.Producer, error)
	GetProducer(uuid string, args Args) (*model.Producer, error)
	GetOneProducerByArgs(args Args) (*model.Producer, error)
	CreateProducer(body *model.Producer) (*model.Producer, error)
	UpdateProducer(uuid string, body *model.Producer, checkAutoMap bool) (*model.Producer, error)
	DeleteProducer(uuid string) (bool, error)
	SyncProducerWriterClones(uuid string) ([]*interfaces.SyncModel, error)
}
type ProducerAPI struct {
	DB ProducerDatabase
}

func (j *ProducerAPI) GetProducers(ctx *gin.Context) {
	args := buildProducerArgs(ctx)
	q, err := j.DB.GetProducers(args)
	ResponseHandler(q, err, ctx)
}

func (j *ProducerAPI) GetProducer(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildProducerArgs(ctx)
	q, err := j.DB.GetProducer(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (j *ProducerAPI) GetOneProducerByArgs(ctx *gin.Context) {
	args := buildProducerArgs(ctx)
	q, err := j.DB.GetOneProducerByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (j *ProducerAPI) CreateProducer(ctx *gin.Context) {
	body, _ := getBODYProducer(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
	}
	body, err = j.DB.CreateProducer(body)
	ResponseHandler(body, err, ctx)
}

func (j *ProducerAPI) UpdateProducer(ctx *gin.Context) {
	body, _ := getBODYProducer(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateProducer(uuid, body, true)
	ResponseHandler(q, err, ctx)
}

func (j *ProducerAPI) DeleteProducer(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteProducer(uuid)
	ResponseHandler(q, err, ctx)
}

func (j *ProducerAPI) SyncProducerWriterClones(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.SyncProducerWriterClones(uuid)
	ResponseHandler(q, err, ctx)
}
