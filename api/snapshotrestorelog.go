package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type SnapshotRestoreLogDatabase interface {
	GetSnapshotRestoreLogs(hostUUID string) ([]*model.SnapshotRestoreLog, error)
	CreateSnapshotRestoreLog(body *model.SnapshotRestoreLog) (*model.SnapshotRestoreLog, error)
	UpdateSnapshotRestoreLog(uuid string, body *model.SnapshotRestoreLog) (*model.SnapshotRestoreLog, error)
	DeleteSnapshotRestoreLog(uuid string) (*interfaces.Message, error)

	ResolveHost(uuid string, name string) (*model.Host, error)
}

type SnapshotRestoreLogAPI struct {
	DB SnapshotRestoreLogDatabase
}

func (a *SnapshotRestoreLogAPI) GetSnapshotRestoreLogs(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q, err := a.DB.GetSnapshotRestoreLogs(host.UUID)
	ResponseHandler(q, err, ctx)
}

func (a *SnapshotRestoreLogAPI) UpdateSnapshotRestoreLog(ctx *gin.Context) {
	body, _ := getBodySnapshotRestoreLog(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateSnapshotRestoreLog(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *SnapshotRestoreLogAPI) DeleteSnapshotRestoreLog(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteSnapshotRestoreLog(uuid)
	ResponseHandler(q, err, ctx)
}
