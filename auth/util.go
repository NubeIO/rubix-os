package auth

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
)

// RegisterAuthentication registers the user id, user and or token.
func RegisterAuthentication(ctx *gin.Context, user *model.User, userID uint, tokenID string) {
	ctx.Set("user", user)
	ctx.Set("userid", userID)
	ctx.Set("tokenid", tokenID)
}

// GetUserID returns the user id which was previously registered by RegisterAuthentication.
func GetUserID(ctx *gin.Context) uint {
	//TODO clean this up
	//user := ctx.MustGet("user").(*model.User)
	//fmt.Println(user, 98989)
	//if user == nil {
	//	userID := ctx.MustGet("userid").(uint)
	//	if userID == 0 {
	//		panic("token and user may not be null")
	//	}
	//	return userID
	//}
	return 1
}

// GetTokenID returns the tokenID.
func GetTokenID(ctx *gin.Context) string {
	return ctx.MustGet("tokenid").(string)
}
