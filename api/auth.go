package api

import (
	"fmt"
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
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
		if len(auth) > 0 {
			// Internal Auth
			if len(auth) == 2 && auth[0] == "Internal" {
				c.Next()
				return
			}
			// Token Auth
			if len(auth) == 2 && auth[0] == "External" {
				fmt.Println(auth[1])
				valid, _ := j.DB.ValidateToken(auth[1])
				if valid {
					c.Next()
					return
				} else {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
			}
			authorized, _ := security.DecodeJwtToken(auth[len(auth)-1], config.Get().SecretKey)
			if authorized {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
