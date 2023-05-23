package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/externaltoken"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"github.com/gin-gonic/gin"
)

// The TokenDatabase interface for encapsulating database access.
type TokenDatabase interface {
}

type TokenAPI struct {
}

func (j *TokenAPI) GetTokens(ctx *gin.Context) {
	q, err := externaltoken.GetExternalTokens()
	ResponseHandler(q, err, ctx)
}

func (j *TokenAPI) GetToken(c *gin.Context) {
	u := c.Param("uuid")
	q, err := externaltoken.GetExternalToken(u)
	ResponseHandler(q, err, c)
}

func (j *TokenAPI) GenerateToken(ctx *gin.Context) {
	body, err := getBodyTokenCreate(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q, err := externaltoken.CreateExternalToken(&externaltoken.ExternalToken{
		UUID:    nuuid.MakeTopicUUID(model.CommonNaming.Token),
		Name:    body.Name,
		Blocked: *body.Blocked})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (j *TokenAPI) RegenerateToken(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := externaltoken.RegenerateExternalToken(uuid)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (j *TokenAPI) BlockToken(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, err := getBodyTokenBlock(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q, err := externaltoken.BlockExternalToken(uuid, *body.Blocked)
	ResponseHandler(q, err, ctx)
}

func (j *TokenAPI) DeleteToken(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := externaltoken.DeleteExternalToken(uuid)
	ResponseHandler(q, err, ctx)
}
