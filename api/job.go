package api

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"strconv"
	"time"
)

var CRON  *cron.Cron

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
	c   *cron.Cron
}
var jobModel *model.Job

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
	ctx.JSON(http.StatusOK, body)

}


func (j *JobAPI) UpdateJob(ctx *gin.Context) {
	body, _ := getBODYJobs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateJob(uuid, body)
	err = j.jobAdd(body)
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
	j.jobRemover(uuid)
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
	res := Data(q)
	ctx.JSON(res.GetStatusCode(), res.GetResponse())
}




func (j *JobAPI) initCron() {
	CRON = cron.New()
	CRON.Start()

}



func (j *JobAPI) jobAdd(body *model.Job) error {
	if body.IsActive {
		//create cron job
		entryID, err := CRON.AddJob(body.Frequency, body)
		if err != nil {
			return err
		}
		body.CronEntryID = int(entryID)

	}
	return nil

}


func (j *JobAPI) jobRemover(id string) {
	value, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		log.Printf("jobRemover %v error \n\n", id)
	}
	CRON.Remove(cron.EntryID(value))
	log.Printf("entry %v is expired and removed \n\n", id)
}


func (j *JobAPI) jobsRemover() {
	for id:=range model.RemoveJob {
		CRON.Remove(cron.EntryID(id))
		log.Printf("entry %v is expired and removed \n",id)
	}
}


func (j *JobAPI) NewJobEngine() {
	j.initCron()
	log.Println("recovering from database if any job exists and not expired")
	job := jobModel
	_jobs, err := j.DB.GetJobs()
	if err != nil {
		log.Fatal(err)
	}
	for _, _job := range _jobs {
		if (_job.EndDate.Unix()<0 || _job.EndDate.UnixNano() > time.Now().UnixNano()) && _job.IsActive {
			entryID, err := CRON.AddJob(_job.Frequency, _job)
			if err != nil {
				log.Println(err)
			}
			log.Println(entryID, job)
		} else if _job.CronEntryID>0 {
			fmt.Println(333333333333333333)
			//job.CronEntryID=0
			//_, err = j.DB.UpdateJob(string(entryID), job)
			//if err != nil {
			//	log.Println(err)
			//}
		}
	}
}

