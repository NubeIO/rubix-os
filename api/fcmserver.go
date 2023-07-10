package api

import (
	"errors"
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
	q.Key = fmt.Sprintf("%s***%s", q.Key[:4], q.Key[len(q.Key)-4:])
	ResponseHandler(q, err, ctx)
}

func (a *FcmServerAPI) UpsertFcmServer(ctx *gin.Context) {
	body, _ := getBodyFcmServer(ctx)
	if len(body.Key) < 20 {
		ResponseHandler(nil, errors.New("key must at least 20 characters in length"), ctx)
		return
	}
	_, err := a.DB.UpsertFcmServer(body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(interfaces.Message{Message: "FCM server has been saved successfully"}, nil, ctx)
}
