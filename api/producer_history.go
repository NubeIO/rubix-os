package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

// The ProducerHistoryDatabase interface for encapsulating database access.
type ProducerHistoryDatabase interface {
	GetProducerHistories(args Args) ([]*model.ProducerHistory, error)
	GetProducerHistoriesByProducerUUID(pUuid string, args Args) ([]*model.ProducerHistory, int64, error)
	GetProducerHistoriesByProducerName(name string) ([]*model.ProducerHistory, int64, error)
	GetLatestProducerHistoryByProducerName(name string) (*model.ProducerHistory, error)
	GetLatestProducerHistoryByProducerUUID(pUuid string) (*model.ProducerHistory, error)
	GetProducerHistoriesByPointUUIDs(pointUUIDs []string, args Args) ([]*interfaces.ProducerHistoryByPointUUID, error)
	GetProducerHistoriesPoints(args Args) ([]*model.History, error)
	GetProducerHistoriesPointsForSync(id string, timeStamp string) ([]*model.History, error)
	DeleteProducerHistoriesByProducerUUID(pUuid string, args Args) (bool, error)
}
type HistoriesAPI struct {
	DB ProducerHistoryDatabase
}

func (a *HistoriesAPI) GetProducerHistories(ctx *gin.Context) {
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.GetProducerHistories(args)
	ResponseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetProducerHistoriesByProducerUUID(ctx *gin.Context) {
	pUuid := resolveProducerUUID(ctx)
	args := buildProducerHistoryArgs(ctx)
	q, _, err := a.DB.GetProducerHistoriesByProducerUUID(pUuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetProducerHistoriesByProducerName(ctx *gin.Context) {
	name := resolveName(ctx)
	q, _, err := a.DB.GetProducerHistoriesByProducerName(name)
	// TODO: @BINOD how do we get the name and count returned with the history data?
	// q, cnt, err := a.DB.GetProducerHistoriesByProducerName(name)
	/*
		type response struct {
			name      string                   `json:"name"`
			count     int64                    `json:"count"`
			histories []*model.ProducerHistory `json:"histories"`
		}
		resp := response{
			name:      name,
			count:     cnt,
			histories: q,
		}
		ResponseHandler(resp, err, ctx)
	*/
	ResponseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetLatestProducerHistoryByProducerName(ctx *gin.Context) {
	name := resolveName(ctx)
	q, err := a.DB.GetLatestProducerHistoryByProducerName(name)
	ResponseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetLatestProducerHistoryByProducerUUID(ctx *gin.Context) {
	pUuid := resolveProducerUUID(ctx)
	q, err := a.DB.GetLatestProducerHistoryByProducerUUID(pUuid)
	ResponseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetProducerHistoriesByPointUUIDs(ctx *gin.Context) {
	args := buildProducerHistoryArgs(ctx)
	var pointUUIDs []string
	err := ctx.ShouldBindJSON(&pointUUIDs)
	q, err := a.DB.GetProducerHistoriesByPointUUIDs(pointUUIDs, args)
	ResponseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetProducerHistoriesPoints(ctx *gin.Context) {
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.GetProducerHistoriesPoints(args)
	ResponseHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetProducerHistoriesPointsForSync(ctx *gin.Context) {
	id, timeStamp := buildProducerHistoryPointsSyncArgs(ctx)
	q, err := a.DB.GetProducerHistoriesPointsForSync(id, timeStamp)
	ResponseHandler(q, err, ctx)
}

func (a *HistoriesAPI) DeleteProducerHistoriesByProducerUUID(ctx *gin.Context) {
	pUuid := resolveProducerUUID(ctx)
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.DeleteProducerHistoriesByProducerUUID(pUuid, args)
	ResponseHandler(q, err, ctx)
}
