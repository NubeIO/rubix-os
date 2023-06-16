package api

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/auth"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// The AuthDatabase interface for encapsulating database access.
type AuthDatabase interface {
	GetMemberSidebars(username string, includeWithoutViews bool) ([]*model.Location, error)
}

type AuthAPI struct {
	DB AuthDatabase
}

func (j *AuthAPI) HandleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorized := auth.Authorize(c.Request)
		if !authorized {
			c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized access"))
			return
		}
		c.Next()
		return
	}
}

func (j *AuthAPI) HandleMemberAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if strings.HasPrefix(authorization, "Internal") || strings.HasPrefix(authorization, "External") {
			if auth.Authorize(c.Request) {
				c.Next()
				return
			}
		}
		username := auth.GetAuthorizedUsername(c.Request)
		if username != "" {
			hostUUID, hostName := matchHostUUIDName(c)
			locations, _ := j.DB.GetMemberSidebars(username, true)
			for _, location := range locations {
				for _, group := range location.Groups {
					for _, host := range group.Hosts {
						if host.UUID == hostUUID || host.Name == hostName {
							c.Next()
							return
						}
					}
				}
			}
		}
		c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized access"))
		return
	}
}
