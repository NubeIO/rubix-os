package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

type PointHistoryDatabase interface {
	GetPointHistories(args Args) ([]*model.PointHistory, error)
	GetPointHistoriesByPointUUID(pUuid string, args Args) ([]*model.PointHistory, int64, error)
	GetLatestPointHistoryByPointUUID(pUuid string) (*model.PointHistory, error)
	GetPointHistoriesByPointUUIDs(pointUUIDs []string, args Args) ([]*model.PointHistory, error)
	GetPointHistoriesForSync(id string, timeStamp string) ([]*model.PointHistory, error)
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
