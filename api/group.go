package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type GroupDatabase interface {
	GetGroups() ([]*model.Group, error)
	GetGroup(uuid string) (*model.Group, error)
	CreateGroup(body *model.Group) (*model.Group, error)
	UpdateGroup(uuid string, body *model.Group) (*model.Group, error)
	DeleteGroup(uuid string) (*interfaces.Message, error)
	DropGroups() (*interfaces.Message, error)
	UpdateHostsStatus(uuid string) (*model.Group, error)
}

type GroupAPI struct {
	DB GroupDatabase
}

func (a *GroupAPI) GetGroupSchema(ctx *gin.Context) {
	q := interfaces.GetGroupSchema()
	ResponseHandler(q, nil, ctx)
}

func (a *GroupAPI) GetGroups(ctx *gin.Context) {
	q, err := a.DB.GetGroups()
	ResponseHandler(q, err, ctx)
}

func (a *GroupAPI) GetGroup(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetGroup(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *GroupAPI) CreateGroup(ctx *gin.Context) {
	body, _ := getBodyGroup(ctx)
	q, err := a.DB.CreateGroup(body)
	ResponseHandler(q, err, ctx)
}

func (a *GroupAPI) UpdateGroup(ctx *gin.Context) {
	body, _ := getBodyGroup(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateGroup(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *GroupAPI) DeleteGroup(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteGroup(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *GroupAPI) DropGroups(ctx *gin.Context) {
	q, err := a.DB.DropGroups()
	ResponseHandler(q, err, ctx)
}

func (a *GroupAPI) UpdateHostsStatus(ctx *gin.Context) {
	hosts, err := a.DB.UpdateHostsStatus(ctx.Params.ByName("uuid"))
	ResponseHandler(hosts, err, ctx)
}
