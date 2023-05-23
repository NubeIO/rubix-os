package api

import (
	"errors"
	"github.com/NubeIO/lib-date/datelib"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/services/system"
	"github.com/NubeIO/rubix-os/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"strconv"
)

var rebootHostJobTag = "reboot.host"

type SystemAPI struct {
	System    *system.System
	Scheduler *gocron.Scheduler
	FileMode  int
}

func (a *SystemAPI) HostTime(c *gin.Context) {
	data := datelib.New(&datelib.Date{}).SystemTime()
	ResponseHandler(data, nil, c)
}

func (a *SystemAPI) GetSystem(c *gin.Context) {
	data, err := a.System.GetSystem()
	ResponseHandler(data, err, c)
}

func (a *SystemAPI) GetMemoryUsage(c *gin.Context) {
	data, err := a.System.GetMemoryUsage()
	ResponseHandler(data, err, c)
}

func (a *SystemAPI) GetMemory(c *gin.Context) {
	data, err := a.System.GetMemory()
	ResponseHandler(data, err, c)
}

func (a *SystemAPI) GetTopProcesses(c *gin.Context) {
	count, err := strconv.Atoi(c.Query("count"))
	m := system.TopProcesses{
		Count: count,
		Sort:  c.Query("sort"),
	}
	data, err := a.System.GetTopProcesses(m)
	ResponseHandler(data, err, c)
}

func (a *SystemAPI) GetSwap(c *gin.Context) {
	data, err := a.System.GetSwap()
	ResponseHandler(data, err, c)
}

func (a *SystemAPI) DiscUsage(c *gin.Context) {
	data, err := a.System.DiscUsage()
	ResponseHandler(data, err, c)
}

func (a *SystemAPI) DiscUsagePretty(c *gin.Context) {
	data, err := a.System.DiscUsagePretty()
	ResponseHandler(data, err, c)
}

func (a *SystemAPI) RunScanner(c *gin.Context) {
	var m *system.Scanner
	err := c.ShouldBindJSON(&m)
	data, err := a.System.RunScanner(m)
	ResponseHandler(data, err, c)
}

func (a *SystemAPI) GetNetworkInterfaces(c *gin.Context) {
	networks, err := nets.GetNetworks()
	ResponseHandler(networks, err, c)
}

func (a *SystemAPI) RebootHost(c *gin.Context) {
	data, err := a.System.RebootHost()
	ResponseHandler(data, err, c)
}

func (a *SystemAPI) GetRebootHostJob(c *gin.Context) {
	rebootJob := utils.GetRebootJob()
	if rebootJob == nil {
		ResponseHandler(nil, errors.New("reboot job not found"), c)
		return
	}
	ResponseHandler(rebootJob, nil, c)
}

func (a *SystemAPI) UpdateRebootHostJob(c *gin.Context) {
	body, err := getBodyRebootJob(c)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	err = utils.ValidateCornExpression(body.Expression)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	body.Tag = rebootHostJobTag
	err = utils.SaveRebootJob(body, a.FileMode)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	_, err = a.Scheduler.Cron(body.Expression).Tag(body.Tag).Do(func() {
		_, _ = a.System.RebootHost()
	})
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(body, nil, c)
}

func (a *SystemAPI) DeleteRebootHostJob(c *gin.Context) {
	err := utils.SaveRebootJob(nil, a.FileMode)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	err = a.Scheduler.RemoveByTag(rebootHostJobTag)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(model.Message{Message: "deleted system reboot job successfully"}, nil, c)
}
