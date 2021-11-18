package api

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The IntegrationDatabase interface for encapsulating database access.
type IntegrationDatabase interface {
	GetIntegration(uuid string) (*model.Integration, error)
	GetIntegrations() ([]*model.Integration, error)
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
	responseHandler(q, err, ctx)
}

func (j *IntegrationAPI) GetIntegrations(ctx *gin.Context) {
	q, err := j.DB.GetIntegrations()
	responseHandler(q, err, ctx)

}

func (j *IntegrationAPI) CreateIntegration(ctx *gin.Context) {
	body, _ := getBODYIntegration(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		responseHandler(nil, err, ctx)
	}
	q, err := j.DB.CreateIntegration(body)
	responseHandler(q, err, ctx)
}

func (j *IntegrationAPI) UpdateIntegration(ctx *gin.Context) {
	body, _ := getBODYIntegration(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateIntegration(uuid, body)
	responseHandler(q, err, ctx)
}

func (j *IntegrationAPI) DeleteIntegration(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteIntegration(uuid)
	responseHandler(q, err, ctx)
}

func (j *IntegrationAPI) DropIntegrationsList(ctx *gin.Context) {
	q, err := j.DB.DropIntegrationsList()
	responseHandler(q, err, ctx)

}
