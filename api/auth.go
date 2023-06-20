package api

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/auth"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/constants"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// The AuthDatabase interface for encapsulating database access.
type AuthDatabase interface {
	GetMemberSidebars(username string, includeWithoutViews bool) ([]*model.Location, error)
}

type AuthAPI struct {
	DB AuthDatabase
}

func (j *AuthAPI) authorizeMember(c *gin.Context) bool {
	username, _ := auth.GetAuthorizedUsername(c.Request)
	if username == "" {
		return false
	}
	hostUUID, hostName := matchHostUUIDName(c)
	locations, _ := j.DB.GetMemberSidebars(username, true)
	for _, location := range locations {
		for _, group := range location.Groups {
			for _, host := range group.Hosts {
				if host.UUID == hostUUID || host.Name == hostName {
					return true
				}
			}
		}
	}
	return false
}

func (j *AuthAPI) HandleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authorized := auth.AuthorizeInternal(c.Request); authorized {
			c.Next()
			return
		}
		if authorized := auth.AuthorizeExternal(c.Request); authorized {
			c.Next()
			return
		}
		if authorized, err := auth.AuthorizeRoles(c.Request, roles...); authorized {
			if utils.Contains(roles, constants.UserRole) {
				c.Next()
				return
			} else {
				if j.authorizeMember(c) {
					c.Next()
					return
				}
			}
		} else if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.AbortWithError(http.StatusUnauthorized, errors.New("token is invalid"))
	}
}
