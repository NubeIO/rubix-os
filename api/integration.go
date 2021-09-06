package api

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The IntegrationDatabase interface for encapsulating database access.
type IntegrationDatabase interface {
	GetIntegration(uuid string) (*model.Integration, error)
	GetIntegrationsList() ([]*model.Integration, error)
	CreateIntegration(body *model.Integration) (*model.Integration, error)
	UpdateIntegration(uuid string, body *model.Integration) (*model.Integration, error)
	DeleteIntegration(uuid string) (bool, error)
	DropIntegrationsList() (bool, error)
}

type IntegrationAPI struct {
	DB IntegrationDatabase
}


func (j *IntegrationAPI) GetIntegration(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetIntegration(uuid)
	reposeHandler(q, err, ctx)
}

func (j *IntegrationAPI) GetIntegrationsList(ctx *gin.Context) {
	q, err := j.DB.GetIntegrationsList()
	reposeHandler(q, err, ctx)

}

func (j *IntegrationAPI) CreateIntegration(ctx *gin.Context) {
	body, _ := getBODYIntegration(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateIntegration(body)
	reposeHandler(q, err, ctx)
}

func (j *IntegrationAPI) UpdateIntegration(ctx *gin.Context) {
	body, _ := getBODYIntegration(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateIntegration(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *IntegrationAPI) DeleteIntegration(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteIntegration(uuid)
	reposeHandler(q, err, ctx)
}


func (j *IntegrationAPI) DropIntegrationsList(ctx *gin.Context) {
	q, err := j.DB.DropIntegrationsList()
	reposeHandler(q, err, ctx)

}
