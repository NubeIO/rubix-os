package auth

import (
	"errors"
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	Authorization = "Authorization"
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
			token := ctx.Request.Header.Get(Authorization)
			if token == "" {
				ctx.AbortWithError(401, errors.New("missing authorization header"))
				return
			}
			if strings.Contains(token, "Internal ") {
				if GetRubixServiceInternalToken() != token {
					ctx.AbortWithError(401, errors.New("internal token mismatch"))
					return
				}
			} else {
				_, err := VerifyToken(token, a.Conf)
				if err != nil {
					ctx.AbortWithError(401, err)
					return
				}
			}
		}
		ctx.Next()
	}
}
