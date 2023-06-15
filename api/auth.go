package api

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/auth"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// The AuthDatabase interface for encapsulating database access.
type AuthDatabase interface {
	GetMemberSidebars(username string) ([]*model.Location, error)
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
		username := auth.GetAuthorizedUsername(c.Request)
		if username != "" {
			hostUUID, hostName := matchHostUUIDName(c)
			locations, _ := j.DB.GetMemberSidebars(username)
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
