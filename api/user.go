package api

import (
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/nerrors"
	"github.com/NubeIO/flow-framework/utils/security"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

// The UserDatabase interface for encapsulating database access.
type UserDatabase interface {
	GetUser() (*model.User, error)
	UpdateUser(body *model.User) (*model.User, error)
}

type UserAPI struct {
	DB UserDatabase
}

func (j *UserAPI) Login(ctx *gin.Context) {
	body, err := getBodyUser(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q, err := j.DB.GetUser()
	if q != nil && body.Username == q.Username && security.CheckPasswordHash(q.Password, body.Password) {
		token, err := security.EncodeJwtToken(q.Username, config.Get().SecretKey)
		if err != nil {
			ResponseHandler(nil, err, ctx)
			return
		}
		ResponseHandler(token, err, ctx)
		return
	}
	ResponseHandler(nil, nerrors.NewErrUnauthorized("check username & password"), ctx)

}

func (j *UserAPI) UpdateUser(ctx *gin.Context) {
	body, err := getBodyUser(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	if !validUsername(body.Username) {
		ResponseHandler("username should be alphanumeric and can contain '_', '-'", nil, ctx)
		return
	}
	body.Password, err = security.GeneratePasswordHash(body.Password)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q, err := j.DB.UpdateUser(body)
	if q != nil {
		q.Password = "******"
	}
	ResponseHandler(q, err, ctx)
}

func (j *UserAPI) GetUser(ctx *gin.Context) {
	q, err := j.DB.GetUser()
	if q != nil {
		q.Password = "******"
	}
	ResponseHandler(q, err, ctx)
}
