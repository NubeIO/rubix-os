package api

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

var CRON *gocron.Scheduler

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

func reposeHandler(body interface{}, err error, ctx *gin.Context) {
	if err != nil {
		if body == nil {
			ctx.JSON(404, "unknown error")
		} else {
			ctx.JSON(404, err.Error())
		}
	} else {
		ctx.JSON(200, body)
	}
}

func (j *JobAPI) GetJobs(ctx *gin.Context) {
	q, err := j.DB.GetJobs()
	reposeHandler(q, err, ctx)

}

func (j *JobAPI) CreateJob(ctx *gin.Context) {
	body, _ := getBODYJobs(ctx)
	_, err := govalidator.ValidateStruct(body)
	q, err := j.DB.CreateJob(body)
	reposeHandler(q, err, ctx)
}

func (j *JobAPI) UpdateJob(ctx *gin.Context) {
	body, _ := getBODYJobs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateJob(uuid, body)
	reposeHandler(q, err, ctx)
}

func (j *JobAPI) GetJob(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetJob(uuid)
	reposeHandler(q, err, ctx)
}

func (j *JobAPI) DeleteJob(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteJob(uuid)
	reposeHandler(q, err, ctx)
}

func (j *JobAPI) initCron() {
	CRON = gocron.NewScheduler(time.UTC)
	CRON.StartAsync()
	//j.syncJobs()
}

//syncJobs start all the jobs on start of the app
//func (j *JobAPI) syncJobs()  {
//	q, err := j.DB.GetJobs()
//	for _, job := range q {
//		if job.Enable {
//			for _, jobSub := range job.JobProducer {
//				if jobSub.Enable {
//					err = j.jobAdd(job.UUID, &job)
//					if err != nil {
//						log.Println("error on read job")
//					}
//				}
//			}
//		}
//	}
//}

//jobAdd add a job
func (j *JobAPI) jobAdd(uuid string, body *model.Job) error {
	if body.Frequency == "" {
		return errors.New("invalid time frequency, example 5m")
	}
	_, err := CRON.Every(body.Frequency).Tag(uuid).Do(taskWithParams, uuid, body)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobAPI) jobRemover(uuid string) error {
	err := CRON.RemoveByTag(uuid)
	if err != nil {
		return errors.New("error on remove job")
	}
	return nil

}

//syncJobs start all the jobs on start of the app
//func (j *JobAPI) syncJobs()  {
//	q, err := j.DB.GetJobs()
//	for _, job := range q {
//		if job.Enable {
//			for _, jobSub := range job.JobProducer {
//				if jobSub.Enable {
//					err = j.jobAdd(job.UUID, &job)
//					if err != nil {
//						log.Println("error on read job")
//					}
//				}
//			}
//		}
//	}
//}

func (j *JobAPI) NewJobEngine() {
	j.initCron()
	log.Info("Init CRON...")
}

func taskWithParams(uuid string, body *model.Job) {
	fmt.Println(uuid)
	//payload := new(payloadBody)
	//payload.UUID = uuid
	//payload.Delete = false
	//payload.MessageString = "what up"
	//payload.MessageTS = time.Now().Format(time.RFC850)
	//topic := fmt.Sprintf("%s:%s", "job",uuid)
	//
	//BUS.RegisterTopics(topic)
	//err := BUS.Emit(BusBackground, topic, payload)
	//
	//fmt.Println("topics", BUS.Topics())
	//if err != nil {
	//	fmt.Println("error", err)
	//}

}
