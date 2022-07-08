package api

import (
	"github.com/NubeIO/flow-framework/nerrors"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/user"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

// The UserDatabase interface for encapsulating database access.
type UserDatabase interface {
}

type UserAPI struct {
}

func (j *UserAPI) Login(ctx *gin.Context) {
	body, err := getBodyUser(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q, err := user.Login(body)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	ResponseHandler(&model.TokenResponse{AccessToken: q, TokenType: "JWT"}, err, ctx)
}

func (j *UserAPI) UpdateUser(ctx *gin.Context) {
	body, err := getBodyUser(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q, err := user.CreateUser(body)
	ResponseHandler(q, err, ctx)
}

func (j *UserAPI) GetUser(ctx *gin.Context) {
	q, err := user.GetUser()
	ResponseHandler(q, err, ctx)
}
