package jobs

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/src/dbhandler"
	"github.com/go-co-op/gocron"
	"time"
)

type Jobs struct {
	db      dbhandler.Handler
	Enabled bool
}

var cron *gocron.Scheduler
var bus eventbus.BusService

func (j *Jobs) InitCron() {
	bus = eventbus.NewService(eventbus.GetBus())
	fmt.Println("IN InitCron")
	cron = gocron.NewScheduler(time.UTC)
	cron.StartAsync()
	j.syncJobs()
	j.Enabled = true

}

func (j *Jobs) task() {
	fmt.Println("TASK RUN")
	t := fmt.Sprintf("%s.%s.%s", eventbus.JobTrigger, "aa", "aa")
	bus.RegisterTopic("job.trigger.aa.aa")
	aa := cron.Jobs()
	for i, e := range aa {
		fmt.Println(i, e.Tags())
	}
	err := bus.Emit(eventbus.CTX(), t, "MESGA OVER BUS")
	if err != nil {

	}

}

func (j *Jobs) syncJobs() {
	fmt.Println("IN JOBS")

	_, err := cron.Cron("*/1 * * * *").Do(j.task)
	if err != nil {
		//return
	}
	//t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsCreated, "aa", "aa")
	//j.bus.RegisterTopic("plugin.created.aa.aa")
	//err = j.bus.Emit(eventbus.CTX(), t, "MESGA OVER BUS")
	//if err != nil {
	//
	//}
	//q, err := j.DB.GetJobs()
	//for _, job := range q {
	//	if c.Jobs == nil {
	//		log.Fatalln("No jobs defined")
	//	}
	//for _, j := range c.Jobs {
	//	if j.Name == "" {
	//		log.Fatalln("Job name not defined")
	//	}
	//	if j.Command == "" {
	//		log.Fatalln("Job command not defined")
	//	}
	//	if j.Frequency == "" {
	//		log.Fatalln("Job frequency not defined")
	//	}
	//}
	//if job.Enable {
	//	for _, jobSub := range job.JobProducer {
	//		if jobSub.Enable {
	//			err = j.jobAdd(job.UUID, &job)
	//			if err != nil {
	//				log.Println("error on read job")
	//			}
	//		}
	//	}
	//}
	//}
}
func (j *Jobs) JobAdd(body *model.Job) error {
	job, err := j.db.CreateJob(body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if job.UUID == "" {
		return errors.New("jobs failed to create a job")
	}
	if job.Frequency == "" {
		return errors.New("invalid time frequency, example 5m")
	}
	_, err = cron.Every(body.Frequency).Tag(job.UUID).Do(j.task)

	if err != nil {
		fmt.Println(job.UUID, err)
		return err
	}
	return nil
}

func (j *Jobs) jobRemover(uuid string) error {
	err := cron.RemoveByTag(uuid)
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
