package api

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

type RestartJobApi struct {
	SystemCtl *systemctl.SystemCtl
	FileMode  int
	Scheduler *gocron.Scheduler
}

func (a *RestartJobApi) GetRestartJob(c *gin.Context) {
	restartJobs := utils.GetRestartJobs()
	if restartJobs == nil {
		ResponseHandler([]interfaces.RestartJob{}, nil, c)
		return
	}
	ResponseHandler(restartJobs, nil, c)
}

func (a *RestartJobApi) UpdateRestartJob(c *gin.Context) {
	body, err := getBodyRestartJob(c)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	err = utils.ValidateCornExpression(body.Expression)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	update := false
	restartJobs := utils.GetRestartJobs()
	for i, restartJob := range restartJobs {
		if restartJob.Unit == body.Unit {
			restartJobs[i] = body
			_ = a.Scheduler.RemoveByTag(restartJob.Unit)
			update = true
			break
		}
	}
	if !update {
		restartJobs = append(restartJobs, body)
	}
	err = utils.SaveRestartJobs(restartJobs, a.FileMode)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	_, err = a.Scheduler.Cron(body.Expression).Tag(body.Unit).Do(func() {
		_ = a.SystemCtl.Restart(body.Unit)
	})
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(body, nil, c)
}

func (a *RestartJobApi) DeleteRestartJob(c *gin.Context) {
	unit := c.Param("unit")
	restartJobs := utils.GetRestartJobs()
	deleted := false
	for i, restartJob := range restartJobs {
		if restartJob.Unit == unit {
			err := a.Scheduler.RemoveByTag(restartJob.Unit)
			if err != nil {
				ResponseHandler(nil, err, c)
				return
			}
			restartJobs = append(restartJobs[:i], restartJobs[i+1:]...)
			deleted = true
			break
		}
	}
	if !deleted {
		err := errors.New("unit not found")
		ResponseHandler(nil, err, c)
		return
	}
	err := utils.SaveRestartJobs(restartJobs, a.FileMode)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(model.Message{Message: fmt.Sprintf("deleted %s restart job successfully", unit)}, nil, c)
}
