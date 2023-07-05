package api

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/auth"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/security"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/user"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/constants"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/nerrors"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"net/http"
)

var invalidMemberTokenError = nerrors.NewErrUnauthorized("invalid member token")
var invalidTokenError = nerrors.NewErrUnauthorized("invalid token")

type MemberDatabase interface {
	GetMembers(args argspkg.Args) ([]*model.Member, error)
	GetMember(uuid string, args argspkg.Args) (*model.Member, error)
	GetMemberByUsername(username string, args argspkg.Args) (*model.Member, error)
	GetMemberByEmail(email string, args argspkg.Args) (*model.Member, error)
	CreateMember(body *model.Member) (*model.Member, error)
	UpdateMember(uuid string, body *model.Member) (*model.Member, error)
	DeleteMember(uuid string) (bool, error)
	DeleteMemberByUsername(username string) (bool, error)
	ChangeMemberPassword(uuid string, password string) (*interfaces.Message, error)
	GetMemberSidebars(username string, includeWithoutViews bool) ([]*model.Location, error)
}

type MemberAPI struct {
	DB MemberDatabase
}

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

func (a *MemberAPI) CreateMember(ctx *gin.Context) {
	body, _ := getBodyMember(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
	}
	q, err := a.DB.CreateMember(body)
	if q != nil {
		q.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) Login(ctx *gin.Context) {
	body, _ := getBodyUser(ctx)
	member, _ := a.DB.GetMemberByUsername(body.Username, argspkg.Args{})
	if member != nil && member.Username == body.Username && security.CheckPasswordHash(member.Password, body.Password) {
		if *member.State != string(model.Verified) {
			ctx.JSON(http.StatusForbidden, interfaces.Message{Message: "member is not verified"})
			return
		}
		token, err := security.EncodeJwtToken(member.Username, constants.MemberRole)
		if err != nil {
			ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
			return
		}
		ResponseHandler(&model.TokenResponse{AccessToken: token, TokenType: "JWT"}, err, ctx)
		return
	}
	ResponseHandler(nil, errors.New("username and password combination is incorrect"), ctx)
}

func (a *MemberAPI) CheckEmail(ctx *gin.Context) {
	email := resolveEmail(ctx)
	q, _ := a.DB.GetMemberByEmail(email, argspkg.Args{})
	message := "email already exists"
	if q == nil {
		message = "email does not exists"
	}
	ResponseHandler(interfaces.Message{Message: message}, nil, ctx)
}

func (a *MemberAPI) CheckUsername(ctx *gin.Context) {
	username := resolveUsername(ctx)
	q, _ := a.DB.GetMemberByUsername(username, argspkg.Args{})
	message := "username already exists"
	if q == nil {
		message = "username does not exists"
	}
	ResponseHandler(interfaces.Message{Message: message}, nil, ctx)
}

func (a *MemberAPI) GetMembers(ctx *gin.Context) {
	args := buildMemberArgs(ctx)
	q, err := a.DB.GetMembers(args)
	for _, m := range q {
		m.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) GetMemberByUUID(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildMemberArgs(ctx)
	q, err := a.DB.GetMember(uuid, args)
	if q != nil {
		q.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) DeleteMemberByUUID(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteMember(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) GetMemberByUsername(ctx *gin.Context) {
	username := resolveUsername(ctx)
	args := buildMemberArgs(ctx)
	q, err := a.DB.GetMemberByUsername(username, args)
	if q != nil {
		q.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) VerifyMember(ctx *gin.Context) {
	username := resolveUsername(ctx)
	member, err := a.DB.GetMemberByUsername(username, argspkg.Args{})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	member.State = nstring.New(string(model.Verified))
	_, err = a.DB.UpdateMember(member.UUID, member)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(interfaces.Message{Message: "member has been verified successfully"}, nil, ctx)
}

func (a *MemberAPI) UpdateMemberByUUID(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBodyMember(ctx)
	q, err := a.DB.UpdateMember(uuid, body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	if q != nil {
		q.MaskPassword()
	}
	ResponseHandler(q, nil, ctx)
}

func (a *MemberAPI) ChangeMemberPassword(ctx *gin.Context) {
	uuid := resolveID(ctx)
	body, _ := getBodyChangePassword(ctx)
	q, err := a.DB.ChangeMemberPassword(uuid, body.NewPassword)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, nil, ctx)
}

func (a *MemberAPI) GetMember(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	args := buildMemberArgs(ctx)
	q, err := a.DB.GetMemberByUsername(username, args)
	if q != nil {
		q.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) UpdateMember(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	member, err := a.DB.GetMemberByUsername(username, argspkg.Args{})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	body, _ := getBodyMember(ctx)
	body.State = member.State
	body.Permission = member.Permission
	q, err := a.DB.UpdateMember(member.UUID, body)
	if q != nil {
		q.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) DeleteMember(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	q, err := a.DB.DeleteMemberByUsername(username)
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) ChangePassword(ctx *gin.Context) {
	body, _ := getBodyChangePassword(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
	}
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	member, err := a.DB.GetMemberByUsername(username, argspkg.Args{})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	if member.Username != body.Username || !security.CheckPasswordHash(member.Password, body.Password) {
		ResponseHandler(nil, nerrors.NewErrUnauthorized("invalid username or password"), ctx)
		return
	}
	q, err := a.DB.ChangeMemberPassword(member.UUID, body.NewPassword)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) RefreshToken(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	token, err := security.EncodeJwtToken(username, constants.MemberRole)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	ResponseHandler(&model.TokenResponse{AccessToken: token, TokenType: "JWT"}, err, ctx)
}

func (a *MemberAPI) GetMemberSidebars(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	q, err := a.DB.GetMemberSidebars(username, false)
	ResponseHandler(q, err, ctx)
}
