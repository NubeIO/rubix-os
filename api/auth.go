package api

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/auth"
	authconstants "github.com/NubeIO/nubeio-rubix-lib-auth-go/constants"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/user"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/constants"
	"github.com/NubeIO/rubix-os/nerrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var invalidMemberTokenError = nerrors.NewErrUnauthorized("invalid member token")
var invalidTokenError = nerrors.NewErrUnauthorized("invalid token")

func getAuthorizedUsername(request *http.Request) (string, error) {
	username, err := auth.GetAuthorizedUsername(request)
	if err != nil {
		return "", nerrors.NewErrUnauthorized(err.Error())
	}
	if username == "" {
		return "", invalidMemberTokenError
	}
	return username, nil
}

func getAuthorizedOrDefaultUsername(request *http.Request) (string, error) {
	if auth.AuthorizeInternal(request) || auth.AuthorizeExternal(request) {
		usr, err := user.GetUser()
		if err != nil {
			return "", err
		}
		return usr.Username, nil
	}
	username, _ := getAuthorizedUsername(request)
	if username != "" {
		return username, nil
	}
	return "", invalidTokenError
}

func getAuthorizedOrDefaultRole(request *http.Request) (string, error) {
	if auth.AuthorizeInternal(request) || auth.AuthorizeExternal(request) {
		return authconstants.UserRole, nil
	}
	role, err := auth.GetAuthorizedRole(request)
	if err != nil {
		return "", nerrors.NewErrUnauthorized(err.Error())
	}
	if role == "" {
		return "", invalidTokenError
	}
	return role, nil
}

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
