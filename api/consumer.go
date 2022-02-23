package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type ConsumersDatabase interface {
	GetConsumers(args Args) ([]*model.Consumer, error)
	GetConsumer(uuid string, args Args) (*model.Consumer, error)
	CreateConsumer(body *model.Consumer) (*model.Consumer, error)
	UpdateConsumer(uuid string, body *model.Consumer) (*model.Consumer, error)
	DeleteConsumer(uuid string) (bool, error)
	DropConsumers() (bool, error)
}
type ConsumersAPI struct {
	DB ConsumersDatabase
}

func (j *ConsumersAPI) GetConsumers(ctx *gin.Context) {
	args := buildConsumerArgs(ctx)
	q, err := j.DB.GetConsumers(args)
	responseHandler(q, err, ctx)
}

func (j *ConsumersAPI) GetConsumer(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildConsumerArgs(ctx)
	q, err := j.DB.GetConsumer(uuid, args)
	responseHandler(q, err, ctx)
}

func (j *ConsumersAPI) CreateConsumer(ctx *gin.Context) {
	body, _ := getBODYConsumer(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		responseHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateConsumer(body)
	responseHandler(q, err, ctx)
}

func (j *ConsumersAPI) UpdateConsumer(ctx *gin.Context) {
	body, _ := getBODYConsumer(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateConsumer(uuid, body)
	responseHandler(q, err, ctx)
}

func (j *ConsumersAPI) DeleteConsumer(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteConsumer(uuid)
	responseHandler(q, err, ctx)
}

func (j *ConsumersAPI) DropConsumers(ctx *gin.Context) {
	q, err := j.DB.DropConsumers()
	responseHandler(q, err, ctx)
}
