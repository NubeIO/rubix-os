package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// The JobDatabase interface for encapsulating database access.
type JobDatabase interface {
	GetJob(uuid string) (*model.Job, error)
	GetJobs() ([]*model.Job, error)
	CreateJob(body *model.Job) (*model.Job, error)
	UpdateJob(uuid string, body *model.Job) (*model.Job, error)
	DeleteJob(uuid string) (bool, error)
}
type JobAPI struct {
	DB JobDatabase
}

func (j *JobAPI) GetJobs(ctx *gin.Context) {
	q, err := j.DB.GetJobs()
	responseHandler(q, err, ctx)
}

func (j *JobAPI) CreateJob(ctx *gin.Context) {
	body, _ := getBODYJobs(ctx)
	_, err := govalidator.ValidateStruct(body)
	q, err := j.DB.CreateJob(body)
	responseHandler(q, err, ctx)
}

func (j *JobAPI) UpdateJob(ctx *gin.Context) {
	body, _ := getBODYJobs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateJob(uuid, body)
	responseHandler(q, err, ctx)
}

func (j *JobAPI) GetJob(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetJob(uuid)
	responseHandler(q, err, ctx)
}

func (j *JobAPI) DeleteJob(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteJob(uuid)
	responseHandler(q, err, ctx)
}
