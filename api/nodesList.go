package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The NodeDatabase interface for encapsulating database access.
type NodeDatabase interface {
	GetNode(uuid string) (*model.Node, error)
	GetNodesList() ([]*model.Node, error)
	CreateNode(body *model.Node) (*model.Node, error)
	UpdateNode(uuid string, body *model.Node) (*model.Node, error)
	DeleteNode(uuid string) (bool, error)
	DropNodesList() (bool, error)
}

type NodeAPI struct {
	DB NodeDatabase
}

func (j *NodeAPI) GetNode(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetNode(uuid)
	reposeHandler(q, err, ctx)
}

func (j *NodeAPI) GetNodesList(ctx *gin.Context) {
	q, err := j.DB.GetNodesList()
	reposeHandler(q, err, ctx)

}

func (j *NodeAPI) CreateNode(ctx *gin.Context) {
	body, _ := getBODYNode(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateNode(body)
	reposeHandler(q, err, ctx)
}

func (j *NodeAPI) UpdateNode(ctx *gin.Context) {
	body, _ := getBODYNode(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateNode(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *NodeAPI) DeleteNode(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteNode(uuid)
	reposeHandler(q, err, ctx)
}


func (j *NodeAPI) DropNodesList(ctx *gin.Context) {
	q, err := j.DB.DropNodesList()
	reposeHandler(q, err, ctx)

}
