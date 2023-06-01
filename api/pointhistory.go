package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type PointHistoryDatabase interface {
	GetPointHistories(args Args) ([]*model.PointHistory, error)
	GetPointHistoriesByPointUUID(pUuid string, args Args) ([]*model.PointHistory, int64, error)
	GetPointHistoriesByPointName(name string) ([]*model.PointHistory, int64, error)
	GetLatestPointHistoryByPointName(name string) (*model.PointHistory, error)
	GetLatestPointHistoryByPointUUID(pUuid string) (*model.PointHistory, error)
	GetPointHistoriesByPointUUIDs(pointUUIDs []string, args Args) ([]*model.PointHistory, error)
	GetPointHistoriesPoints(args Args) ([]*model.History, error)
	GetPointHistoriesForSync(id string, timeStamp string) ([]*model.History, error)
	DeletePointHistoriesByPointUUID(pUuid string, args Args) (bool, error)
}
type PointHistoryAPI struct {
	DB PointHistoryDatabase
}

func (a *PointHistoryAPI) GetPointHistories(ctx *gin.Context) {
	args := buildPointHistoryArgs(ctx)
	q, err := a.DB.GetPointHistories(args)
	ResponseHandler(q, err, ctx)
}

func (a *PointHistoryAPI) GetPointHistoriesByPointUUID(ctx *gin.Context) {
	pUuid := resolvePointUUID(ctx)
	args := buildPointHistoryArgs(ctx)
	q, _, err := a.DB.GetPointHistoriesByPointUUID(pUuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *PointHistoryAPI) GetPointHistoriesByPointName(ctx *gin.Context) {
	name := resolveName(ctx)
	q, _, err := a.DB.GetPointHistoriesByPointName(name)
	// TODO: @BINOD how do we get the name and count returned with the history data?
	// q, cnt, err := a.DB.GetPointHistoriesByPointName(name)
	/*
		type response struct {
			name      string                   `json:"name"`
			count     int64                    `json:"count"`
			histories []*model.PointHistory    `json:"histories"`
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

func (a *PointHistoryAPI) GetLatestPointHistoryByPointName(ctx *gin.Context) {
	name := resolveName(ctx)
	q, err := a.DB.GetLatestPointHistoryByPointName(name)
	ResponseHandler(q, err, ctx)
}

func (a *PointHistoryAPI) GetLatestPointHistoryByPointUUID(ctx *gin.Context) {
	pUuid := resolvePointUUID(ctx)
	q, err := a.DB.GetLatestPointHistoryByPointUUID(pUuid)
	ResponseHandler(q, err, ctx)
}

func (a *PointHistoryAPI) GetPointHistoriesByPointUUIDs(ctx *gin.Context) {
	args := buildPointHistoryArgs(ctx)
	var pointUUIDs []string
	err := ctx.ShouldBindJSON(&pointUUIDs)
	q, err := a.DB.GetPointHistoriesByPointUUIDs(pointUUIDs, args)
	ResponseHandler(q, err, ctx)
}

func (a *PointHistoryAPI) GetPointHistoriesPoints(ctx *gin.Context) {
	args := buildPointHistoryArgs(ctx)
	q, err := a.DB.GetPointHistoriesPoints(args)
	ResponseHandler(q, err, ctx)
}

func (a *PointHistoryAPI) GetPointHistoriesForSync(ctx *gin.Context) {
	id, timeStamp := buildPointHistorySyncArgs(ctx)
	q, err := a.DB.GetPointHistoriesForSync(id, timeStamp)
	ResponseHandler(q, err, ctx)
}

func (a *PointHistoryAPI) DeletePointHistoriesByPointUUID(ctx *gin.Context) {
	pUuid := resolvePointUUID(ctx)
	args := buildPointHistoryArgs(ctx)
	q, err := a.DB.DeletePointHistoriesByPointUUID(pUuid, args)
	ResponseHandler(q, err, ctx)
}
