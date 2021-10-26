package api

import (
	"github.com/NubeDev/flow-framework/auth"
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

type LoginAPI struct {
	Conf *config.Configuration
}

func (a *LoginAPI) Login(ctx *gin.Context) {
	body, _ := getBodyLogin(ctx)
	token, err := auth.GetToken(body.Username, body.Password, a.Conf)
	outToken := model.Token{}
	if token != nil {
		outToken = model.Token{
			AccessToken: *token,
			TokenType:   "JWT",
		}
	}
	reposeHandler(outToken, err, ctx)
}
