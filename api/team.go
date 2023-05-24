package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type TeamDatabase interface {
	GetTeams() ([]*model.Team, error)
	GetTeam(uuid string) (*model.Team, error)
	CreateTeam(body *model.Team) (*model.Team, error)
	UpdateTeam(uuid string, body *model.Team) (*model.Team, error)
	DeleteTeam(uuid string) (bool, error)
	DropTeams() (bool, error)
}

type TeamAPI struct {
	DB TeamDatabase
}

func (a *TeamAPI) GetTeams(ctx *gin.Context) {
	q, err := a.DB.GetTeams()
	ResponseHandler(q, err, ctx)
}

func (a *TeamAPI) GetTeam(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetTeam(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *TeamAPI) CreateTeam(ctx *gin.Context) {
	body, _ := getBodyTeam(ctx)
	q, err := a.DB.CreateTeam(body)
	ResponseHandler(q, err, ctx)
}

func (a *TeamAPI) UpdateTeam(ctx *gin.Context) {
	body, _ := getBodyTeam(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateTeam(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *TeamAPI) DeleteTeam(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteTeam(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *TeamAPI) DropTeams(ctx *gin.Context) {
	q, err := a.DB.DropTeams()
	ResponseHandler(q, err, ctx)
}
