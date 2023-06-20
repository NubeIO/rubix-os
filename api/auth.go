package api

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/auth"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/constants"
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

func (j *AuthAPI) authorizeMember(c *gin.Context) (int, error) {
	username, _ := auth.GetAuthorizedUsername(c.Request)
	if username == "" {
		return http.StatusUnauthorized, invalidMemberTokenError
	}
	hostUUID, hostName := matchHostUUIDName(c)
	locations, _ := j.DB.GetMemberSidebars(username, true)
	for _, location := range locations {
		for _, group := range location.Groups {
			for _, host := range group.Hosts {
				if host.UUID == hostUUID || host.Name == hostName {
					return 0, nil
				}
			}
		}
	}
	return http.StatusForbidden, errors.New("forbidden access for the member")
}

func (j *AuthAPI) HandleAuth(hostLevel bool, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authorized := auth.AuthorizeInternal(c.Request); authorized {
			c.Next()
			return
		}
		if authorized := auth.AuthorizeExternal(c.Request); authorized {
			c.Next()
			return
		}
		if authorized, role, err := auth.AuthorizeRoles(c.Request, roles...); authorized {
			if hostLevel && *role == constants.MemberRole {
				if statusCode, err := j.authorizeMember(c); err != nil {
					c.AbortWithError(statusCode, err)
					return
				}
			}
			c.Next()
			return
		} else if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.AbortWithError(http.StatusUnauthorized, errors.New("token is invalid"))
	}
}
