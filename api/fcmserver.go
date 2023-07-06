package api

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type FcmServerDatabase interface {
	GetFcmServer() (*model.FcmServer, error)
	UpsertFcmServer(body *model.FcmServer) (*model.FcmServer, error)
}

type FcmServerAPI struct {
	DB FcmServerDatabase
}

func (a *FcmServerAPI) GetFcmServer(ctx *gin.Context) {
	q, err := a.DB.GetFcmServer()
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q.Key = fmt.Sprintf("%s***%s", q.Key[:1], q.Key[len(q.Key)-1:])
	ResponseHandler(q, err, ctx)
}

func (a *FcmServerAPI) UpsertFcmServer(ctx *gin.Context) {
	body, _ := getBodyFcmServer(ctx)
	_, err := a.DB.UpsertFcmServer(body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(interfaces.Message{Message: "fcm server has been saved successfully"}, nil, ctx)
}
