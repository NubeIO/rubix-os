package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The NodeDatabase interface for encapsulating database access.
type NodeDatabase interface {
	GetNodeList(uuid string) (*model.NodeList, error)
	GetNodesList() ([]*model.NodeList, error)
	CreateNodeList(body *model.NodeList) (*model.NodeList, error)
	UpdateNodeList(uuid string, body *model.NodeList) (*model.NodeList, error)
	DeleteNodeList(uuid string) (bool, error)
	DropNodesList() (bool, error)
}

type NodeListAPI struct {
	DB NodeDatabase
}

func (j *NodeListAPI) GetNodeList(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetNodeList(uuid)
	reposeHandler(q, err, ctx)
}

func (j *NodeListAPI) GetNodesList(ctx *gin.Context) {
	q, err := j.DB.GetNodesList()
	reposeHandler(q, err, ctx)

}

func (j *NodeListAPI) CreateNodeList(ctx *gin.Context) {
	body, _ := getBODYNode(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateNodeList(body)
	reposeHandler(q, err, ctx)
}

func (j *NodeListAPI) UpdateNode(ctx *gin.Context) {
	body, _ := getBODYNode(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateNodeList(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *NodeListAPI) DeleteNodeList(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteNodeList(uuid)
	reposeHandler(q, err, ctx)
}


func (j *NodeListAPI) DropNodesList(ctx *gin.Context) {
	q, err := j.DB.DropNodesList()
	reposeHandler(q, err, ctx)

}
