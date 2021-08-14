package test

import (
	"github.com/NubeDev/plug-framework/model"
	"github.com/gin-gonic/gin"
)

// WithUser fake an authentication for testing.
func WithUser(ctx *gin.Context, userID uint) {
	ctx.Set("user", &model.User{ID: userID})
	ctx.Set("userid", userID)
}
