package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type SnapshotCreteLogDatabase interface {
	GetSnapshotCreateLogs(hostUUID string) ([]*model.SnapshotCreateLog, error)
	CreateSnapshotCreateLog(body *model.SnapshotCreateLog) (*model.SnapshotCreateLog, error)
	UpdateSnapshotCreateLog(uuid string, body *model.SnapshotCreateLog) (*model.SnapshotCreateLog, error)
	DeleteSnapshotCreateLog(uuid string) (*interfaces.Message, error)

	ResolveHost(uuid string, name string) (*model.Host, error)
}

type SnapshotCreateLogAPI struct {
	DB SnapshotCreteLogDatabase
}

func (a *SnapshotCreateLogAPI) GetSnapshotCreateLogs(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q, err := a.DB.GetSnapshotCreateLogs(host.UUID)
	ResponseHandler(q, err, ctx)
}

func (a *SnapshotCreateLogAPI) UpdateSnapshotCreateLog(ctx *gin.Context) {
	body, _ := getBodySnapshotCreateLog(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateSnapshotCreateLog(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *SnapshotCreateLogAPI) DeleteSnapshotCreateLog(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteSnapshotCreateLog(uuid)
	ResponseHandler(q, err, ctx)
}
