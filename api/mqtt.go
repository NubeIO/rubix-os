package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The MqttConnectionDatabase interface for encapsulating database access.
type MqttConnectionDatabase interface {
	GetMqttConnection(uuid string) (*model.MqttConnection, error)
	GetMqttConnectionsList() ([]*model.MqttConnection, error)
	CreateMqttConnection(body *model.MqttConnection) (*model.MqttConnection, error)
	UpdateMqttConnection(uuid string, body *model.MqttConnection) (*model.MqttConnection, error)
	DeleteMqttConnection(uuid string) (bool, error)
	DropMqttConnectionsList() (bool, error)
}

type MqttConnectionAPI struct {
	DB MqttConnectionDatabase
}

func (j *MqttConnectionAPI) GetMqttConnection(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetMqttConnection(uuid)
	reposeHandler(q, err, ctx)
}

func (j *MqttConnectionAPI) GetMqttConnectionsList(ctx *gin.Context) {
	q, err := j.DB.GetMqttConnectionsList()
	reposeHandler(q, err, ctx)

}

func (j *MqttConnectionAPI) CreateMqttConnection(ctx *gin.Context) {
	body, _ := getBODYMqttConnection(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateMqttConnection(body)
	reposeHandler(q, err, ctx)
}

func (j *MqttConnectionAPI) UpdateMqttConnection(ctx *gin.Context) {
	body, _ := getBODYMqttConnection(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateMqttConnection(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *MqttConnectionAPI) DeleteMqttConnection(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteMqttConnection(uuid)
	reposeHandler(q, err, ctx)
}


func (j *MqttConnectionAPI) DropMqttConnectionsList(ctx *gin.Context) {
	q, err := j.DB.DropMqttConnectionsList()
	reposeHandler(q, err, ctx)

}
