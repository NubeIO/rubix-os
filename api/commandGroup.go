package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

/*
Stream
*/

// The CommandGroupDatabase interface for encapsulating database access.
type CommandGroupDatabase interface {
	GetCommandGroup(uuid string) (*model.CommandGroup, error)
	GetCommandGroups() ([]*model.CommandGroup, error)
	CreateCommandGroup(body *model.CommandGroup) (*model.CommandGroup, error)
	UpdateCommandGroup(uuid string, body *model.CommandGroup) (*model.CommandGroup, error)
	DeleteCommandGroup(uuid string) (bool, error)
}

type CommandGroupAPI struct {
	DB CommandGroupDatabase
}

func (j *CommandGroupAPI) GetCommandGroup(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetCommandGroup(uuid)
	reposeHandler(q, err, ctx)
}

func (j *CommandGroupAPI) GetCommandGroups(ctx *gin.Context) {
	q, err := j.DB.GetCommandGroups()
	reposeHandler(q, err, ctx)

}

func (j *CommandGroupAPI) CreateCommandGroup(ctx *gin.Context) {
	body, _ := getBODYCommandGroup(ctx)
	_, err := govalidator.ValidateStruct(body)
	q, err := j.DB.CreateCommandGroup(body)
	reposeHandler(q, err, ctx)
}

func (j *CommandGroupAPI) UpdateCommandGroup(ctx *gin.Context) {
	body, _ := getBODYCommandGroup(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateCommandGroup(uuid, body)
	reposeHandler(q, err, ctx)
}

func (j *CommandGroupAPI) DeleteCommandGroup(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteCommandGroup(uuid)
	reposeHandler(q, err, ctx)
}
