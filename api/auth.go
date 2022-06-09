package api

import (
	"errors"
	"github.com/NubeIO/flow-framework/auth"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/utils/security"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// The AuthDatabase interface for encapsulating database access.
type AuthDatabase interface {
	ValidateToken(accessToken string) (bool, error)
}

type AuthAPI struct {
	DB AuthDatabase
}

func (j *AuthAPI) HandleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(authorization) > 0 {
			// Internal Auth
			if len(authorization) == 2 && authorization[0] == "Internal" &&
				authorization[1] == auth.GetRubixServiceInternalToken(false) {
				c.Next()
				return
			}
			// Token Auth
			if len(authorization) == 2 && authorization[0] == "External" {
				valid, _ := j.DB.ValidateToken(authorization[1])
				if valid {
					c.Next()
					return
				} else {
					c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized access"))
					return
				}
			}
			authorized, _ := security.DecodeJwtToken(authorization[len(authorization)-1], config.Get().SecretKey)
			if authorized {
				c.Next()
				return
			}
		}
		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized access"))
		return
	}
}
