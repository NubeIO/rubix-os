package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	parentArgs "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type LocationDatabase interface {
	GetLocations(args parentArgs.Args) ([]*model.Location, error)
	GetLocation(uuid string, args parentArgs.Args) (*model.Location, error)
	CreateLocation(body *model.Location) (*model.Location, error)
	UpdateLocation(uuid string, body *model.Location) (*model.Location, error)
	DeleteLocation(uuid string) (*interfaces.Message, error)
	DropLocations() (*interfaces.Message, error)
}

type LocationAPI struct {
	DB LocationDatabase
}

func (a *LocationAPI) GetLocationSchema(ctx *gin.Context) {
	q := interfaces.GetLocationSchema()
	ResponseHandler(q, nil, ctx)
}

func (a *LocationAPI) GetLocations(ctx *gin.Context) {
	args := buildLocationArgs(ctx)
	q, err := a.DB.GetLocations(args)
	ResponseHandler(q, err, ctx)
}

func (a *LocationAPI) GetLocation(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildLocationArgs(ctx)
	q, err := a.DB.GetLocation(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *LocationAPI) CreateLocation(ctx *gin.Context) {
	body, _ := getBodyLocation(ctx)
	q, err := a.DB.CreateLocation(body)
	ResponseHandler(q, err, ctx)
}

func (a *LocationAPI) UpdateLocation(ctx *gin.Context) {
	body, _ := getBodyLocation(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateLocation(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *LocationAPI) DeleteLocation(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteLocation(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *LocationAPI) DropLocations(ctx *gin.Context) {
	q, err := a.DB.DropLocations()
	ResponseHandler(q, err, ctx)
}
