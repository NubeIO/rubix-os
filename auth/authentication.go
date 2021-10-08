package auth

import (
	"errors"
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationToken = "Authorization"
)

type Database interface {
	GetApplicationByToken(token string) (*model.Application, error)
	GetClientByToken(token string) (*model.Client, error)
	GetPluginConfByToken(token string) (*model.PluginConf, error)
	GetUserByName(name string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
}

type Auth struct {
	Conf *config.Configuration
}

func (a *Auth) RequireValidToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if a.Conf.Prod || a.Conf.Auth {
			token := a.tokenFromQueryOrHeader(ctx)
			if token == "" {
				ctx.AbortWithError(401, errors.New("missing authorization header"))
				return
			}
			_, err := VerifyToken(token, a.Conf)
			if err != nil {
				ctx.AbortWithError(401, err)
				return
			}
		}
		ctx.Next()
	}
}

func (a *Auth) tokenFromQueryOrHeader(ctx *gin.Context) string {
	if token := a.tokenFromQuery(ctx); token != "" {
		return token
	} else if token = a.tokenFromHeader(ctx); token != "" {
		return token
	}
	return ""
}

func (a *Auth) tokenFromQuery(ctx *gin.Context) string {
	return ctx.Request.URL.Query().Get("token")
}

func (a *Auth) tokenFromHeader(ctx *gin.Context) string {
	if ctx.Request.Header.Get(AuthorizationToken) != "" {
		return ctx.Request.Header.Get(AuthorizationToken)
	}
	return ""
}
