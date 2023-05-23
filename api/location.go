package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type LocationDatabase interface {
	GetLocations() ([]*model.Location, error)
	GetLocation(uuid string) (*model.Location, error)
	CreateLocation(body *model.Location) (*model.Location, error)
	UpdateLocation(uuid string, body *model.Location) (*model.Location, error)
	DeleteLocation(uuid string) (bool, error)
	DropLocations() (bool, error)
}

type LocationAPI struct {
	DB LocationDatabase
}

func (a *LocationAPI) GetLocationSchema(ctx *gin.Context) {
	q := interfaces.GetLocationSchema()
	ResponseHandler(q, nil, ctx)
}

func (a *LocationAPI) GetLocations(ctx *gin.Context) {
	q, err := a.DB.GetLocations()
	ResponseHandler(q, err, ctx)
}

func (a *LocationAPI) GetLocation(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetLocation(uuid)
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
