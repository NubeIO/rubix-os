package api

import (
	"errors"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type StreamLogDatabase interface {
	CreateLogAndReturn(body *interfaces.StreamLog) (*interfaces.StreamLog, error)
	GetStreamsLogs() []*interfaces.StreamLog
	GetStreamLog(uuid string) *interfaces.StreamLog
	CreateStreamLog(body *interfaces.StreamLog) (string, error)
	DeleteStreamLog(uuid string) bool
	DeleteStreamLogs()
}
type StreamLogApi struct {
	DB StreamLogDatabase
}

func (a *StreamLogApi) GetStreamLogs(c *gin.Context) {
	q := a.DB.GetStreamsLogs()
	ResponseHandler(q, nil, c)
}

func (a *StreamLogApi) GetStreamLog(c *gin.Context) {
	u := c.Param("uuid")
	logStream := a.DB.GetStreamLog(u)
	if logStream == nil {
		ResponseHandler(nil, errors.New("log not found"), c)
		return
	}
	ResponseHandler(logStream, nil, c)
}

func (a *StreamLogApi) CreateStreamLog(c *gin.Context) {
	body, err := getBodyStreamLog(c)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	uuid, err := a.DB.CreateStreamLog(body)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(map[string]interface{}{"UUID": uuid}, err, c)
}

func (a *StreamLogApi) CreateLogAndReturn(c *gin.Context) {
	body, err := getBodyStreamLog(c)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	logStream, err := a.DB.CreateLogAndReturn(body)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(logStream, nil, c)
}

func (a *StreamLogApi) DeleteStreamLog(c *gin.Context) {
	uuid := resolveID(c)
	deleted := a.DB.DeleteStreamLog(uuid)
	if !deleted {
		ResponseHandler(nil, errors.New("log not found"), c)
		return
	}
	ResponseHandler(deleted, nil, c)
}

func (a *StreamLogApi) DeleteStreamLogs(c *gin.Context) {
	a.DB.DeleteStreamLogs()
	ResponseHandler(true, nil, c)
}
