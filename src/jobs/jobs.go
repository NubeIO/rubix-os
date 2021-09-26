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
	cron = gocron.NewScheduler(time.UTC)
	cron.StartAsync()
	j.Enabled = true

}

func (j *Jobs) task() {
	fmt.Println("TASK RUN")
	t := fmt.Sprintf("%s.%s.%s", eventbus.JobTrigger, "aa", "aa")
	bus.RegisterTopic("job.trigger.aa.aa") //TODO send via job uuid
	//aa := cron.Jobs()
	//for i, e := range aa {
	//	//fmt.Println("JOB TAGS", i, e.Tags())
	//}
	err := bus.Emit(eventbus.CTX(), t, "MEG OVER BUS") //TODO send via job uuid
	if err != nil {
		//TODO FIX ERROR
	}
}

func (j *Jobs) JobAdd(body *model.Job) error {
	job, err := j.db.CreateJob(body) //TODO is being added to the DB for the plugin influx but maybe its not needed
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
