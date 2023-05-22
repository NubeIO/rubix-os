package api

import (
	"github.com/NubeIO/flow-framework/services/system"
	"github.com/NubeIO/lib-date/datectl"
	"github.com/gin-gonic/gin"
)

type DateAPI struct {
	System *system.System
}

func (a *DateAPI) SystemTime(c *gin.Context) {
	data := a.System.SystemTime()
	ResponseHandler(data, nil, c)
}

func (a *DateAPI) GenerateTimeSyncConfig(c *gin.Context) {
	var m *datectl.TimeSyncConfig
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data := a.System.GenerateTimeSyncConfig(m)
	ResponseHandler(data, nil, c)
}

func (a *DateAPI) GetHardwareTZ(c *gin.Context) {
	data, err := a.System.GetHardwareTZ()
	ResponseHandler(data, err, c)
}

func (a *DateAPI) GetTimeZoneList(c *gin.Context) {
	data, err := a.System.GetTimeZoneList()
	ResponseHandler(data, err, c)
}

func (a *DateAPI) UpdateTimezone(c *gin.Context) {
	var m system.DateBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.UpdateTimezone(m)
	ResponseHandler(data, err, c)
}

func (a *DateAPI) SetSystemTime(c *gin.Context) {
	var m system.DateBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	data, err := a.System.SetSystemTime(m)
	ResponseHandler(data, err, c)
}

func (a *DateAPI) NTPEnable(c *gin.Context) {
	data, err := a.System.NTPEnable()
	ResponseHandler(data, err, c)
}

func (a *DateAPI) NTPDisable(c *gin.Context) {
	data, err := a.System.NTPDisable()
	ResponseHandler(data, err, c)
}
