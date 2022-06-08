package api

import (
	"github.com/NubeIO/flow-framework/utils/security"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
)

// The TokenDatabase interface for encapsulating database access.
type TokenDatabase interface {
	GetTokens() ([]*model.Token, error)
	CreateToken(*model.Token) (*model.Token, error)
	UpdateToken(*model.Token) (*model.Token, error)
}

type TokenAPI struct {
	DB TokenDatabase
}

func (j *TokenAPI) GetTokens(ctx *gin.Context) {
	q, err := j.DB.GetTokens()
	ResponseHandler(q, err, ctx)
}

func (j *TokenAPI) GenerateToken(ctx *gin.Context) {
	body, err := getBodyToken(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	token := &model.Token{Name: body.Name, Token: security.GenerateToken(), Blocked: body.Blocked}
	q, err := j.DB.CreateToken(token)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (j *TokenAPI) UpdateToken(ctx *gin.Context) {
	body, err := getBodyToken(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	token := &model.Token{Name: body.Name, Token: "", Blocked: body.Blocked}
	q, err := j.DB.UpdateToken(token)
	ResponseHandler(q, err, ctx)
}
