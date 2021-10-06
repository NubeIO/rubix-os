package api

import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// The ProducerHistoryDatabase interface for encapsulating database access.
type ProducerHistoryDatabase interface {
	GetProducerHistories(args Args) ([]*model.ProducerHistory, error)
	GetProducerHistoriesByProducerUUID(pUuid string, args Args) ([]*model.ProducerHistory, int64, error)
	GetLatestProducerHistoryByProducerUUID(pUuid string) (*model.ProducerHistory, error)
	CreateBulkProducerHistory(histories []*model.ProducerHistory) (bool, error)
	CreateProducerHistory(history *model.ProducerHistory) (bool, error)
	DeleteAllProducerHistories(args Args) (bool, error)
	DeleteProducerHistoriesByProducerUUID(pUuid string, args Args) (bool, error)
}
type HistoriesAPI struct {
	DB ProducerHistoryDatabase
}

func (a *HistoriesAPI) GetProducerHistories(ctx *gin.Context) {
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.GetProducerHistories(args)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetProducerHistoriesByProducerUUID(ctx *gin.Context) {
	pUuid := resolveProducerUUID(ctx)
	args := buildProducerHistoryArgs(ctx)
	q, _, err := a.DB.GetProducerHistoriesByProducerUUID(pUuid, args)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) GetLatestProducerHistoryByProducerUUID(ctx *gin.Context) {
	pUuid := resolveProducerUUID(ctx)
	fmt.Println(pUuid)
	q, err := a.DB.GetLatestProducerHistoryByProducerUUID(pUuid)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) CreateProducerHistory(ctx *gin.Context) {
	body, _ := getBodyHistory(ctx)
	q, err := a.DB.CreateProducerHistory(body)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) CreateBulkProducerHistory(ctx *gin.Context) {
	body, _ := getBodyBulkHistory(ctx)
	q, err := a.DB.CreateBulkProducerHistory(body)
	reposeHandler(q, err, ctx)
}

func (a *HistoriesAPI) DeleteAllProducerHistories(ctx *gin.Context) {
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.DeleteAllProducerHistories(args)
	reposeHandler(q, err, ctx)

}

func (a *HistoriesAPI) DeleteProducerHistoriesByProducerUUID(ctx *gin.Context) {
	pUuid := resolveProducerUUID(ctx)
	args := buildProducerHistoryArgs(ctx)
	q, err := a.DB.DeleteProducerHistoriesByProducerUUID(pUuid, args)
	reposeHandler(q, err, ctx)
}
