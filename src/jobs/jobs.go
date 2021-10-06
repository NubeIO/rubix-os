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
	db         dbhandler.Handler
	Enabled    bool
}

var cron *gocron.Scheduler
var bus eventbus.BusService

func (j *Jobs) InitCron() {
	bus = eventbus.NewService(eventbus.GetBus())
	cron = gocron.NewScheduler(time.UTC)
	cron.StartAsync()
	j.Enabled = true

}

func (j *Jobs) task(mp string, uuid string) {
	fmt.Println("TASK RUN")
	t := fmt.Sprintf("%s.%s.%s", eventbus.JobTrigger, mp, uuid)
	bus.RegisterTopic(t)
	err := bus.Emit(eventbus.CTX(), t, "MEG OVER BUS")
	if err != nil {
		//TODO FIX ERROR
	}
}

func (j *Jobs) JobAdd(body *model.Job) error {
	pc, err := j.db.GetPlugin(body.PluginConfId)
	if err != nil {
		return err
	}
	_, err = cron.Every(body.Frequency).Tag(body.UUID).Do(j.task, pc.ModulePath, body.UUID)
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
