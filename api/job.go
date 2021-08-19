package api

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"log"
	"net/http"
	"time"
)

var CRON  *gocron.Scheduler

// The JobDatabase interface for encapsulating database access.
type JobDatabase interface {
	GetJob(uuid string) (*model.Job, error)
	GetJobs() ([]model.Job, error)
	CreateJob(body *model.Job) error
	UpdateJob(uuid string, body *model.Job) (*model.Job, error)
	DeleteJob(uuid string) (bool, error)

}
type JobAPI struct {
	DB JobDatabase
}

func (j *JobAPI) GetJobs(ctx *gin.Context) {
	q, err := j.DB.GetJobs()
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())
}

func (j *JobAPI) CreateJob(ctx *gin.Context) {
	body, _ := getBODYJobs(ctx)
	if success := successOrAbort(ctx, http.StatusInternalServerError, j.DB.CreateJob(body)); !success {
		return
	}
	err := j.jobAdd(body.UUID, body)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	ctx.JSON(http.StatusOK, body)

}


func (j *JobAPI) UpdateJob(ctx *gin.Context) {
	body, _ := getBODYJobs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateJob(uuid, body)
	err = j.jobAdd(uuid, body)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}

func (j *JobAPI) GetJob(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetJob(uuid)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())
}


func (j *JobAPI) DeleteJob(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteJob(uuid)
	err = j.jobRemover(uuid)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())
}

func (j *JobAPI) initCron() {
	CRON = 	gocron.NewScheduler(time.UTC)
	CRON.StartAsync()

}

func (j *JobAPI) jobAdd(uuid string, body *model.Job) error {
	if body.Frequency == "" {
		return errors.New("invalid time frequency, example 5m")
	}
	if body.IsActive {
		_, err := CRON.Every(body.Frequency).Tag(uuid).Do(taskWithParams, uuid, "hello")
		if err != nil {
			return err
		}
	}
	return nil

}

func taskWithParams(a string, b string) {
	fmt.Println(time.Now().Format(time.RFC850))
	fmt.Println(a, b)
}

func (j *JobAPI) jobRemover(uuid string) error {
	err := CRON.RemoveByTag(uuid)
	if err != nil {
		log.Printf("job %v was NOT removed \n\n", uuid)
		return errors.New("error on remove job")
	}
	log.Printf("job %v was removed \n\n", uuid)
	return nil

}

func (j *JobAPI) NewJobEngine() {
	j.initCron()
	log.Println("INIT CRON")
}


