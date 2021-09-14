package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The ConsumersDatabase interface for encapsulating database access.
type ConsumersDatabase interface {
	GetConsumer(uuid string, withChildren bool) (*model.Consumer, error)
	GetConsumers() ([]*model.Consumer, error)
	CreateConsumer(body *model.Consumer) (*model.Consumer, error)
	UpdateConsumer(uuid string, body *model.Consumer) (*model.Consumer, error)
	DeleteConsumer(uuid string) (bool, error)
	DropConsumers() (bool, error)
}
type ConsumersAPI struct {
	DB ConsumersDatabase
}

func (j *ConsumersAPI) GetConsumer(ctx *gin.Context) {
	uuid := resolveID(ctx)
	withChildren, _, _ := withChildrenArgs(ctx)
	q, err := j.DB.GetConsumer(uuid, withChildren)
	reposeHandler(q, err, ctx)
}

func (j *ConsumersAPI) GetConsumers(ctx *gin.Context) {
	q, err := j.DB.GetConsumers()
	reposeHandler(q, err, ctx)
}

func (j *ConsumersAPI) CreateConsumer(ctx *gin.Context) {
	body, _ := getBODYConsumer(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateConsumer(body)
	reposeHandler(q, err, ctx)
}

func (j *ConsumersAPI) UpdateConsumer(ctx *gin.Context) {
	body, _ := getBODYConsumer(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateConsumer(uuid, body)
	reposeHandler(q, err, ctx)
}

func (j *ConsumersAPI) DeleteConsumer(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteConsumer(uuid)
	reposeHandler(q, err, ctx)
}

func (j *ConsumersAPI) DropConsumers(ctx *gin.Context) {
	q, err := j.DB.DropConsumers()
	reposeHandler(q, err, ctx)
}
