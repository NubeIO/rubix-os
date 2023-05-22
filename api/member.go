package api

import (
	"errors"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/nerrors"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/auth"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/utils/security"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

var invalidMemberTokenError = nerrors.NewErrUnauthorized("invalid member token")

type MemberDatabase interface {
	GetMembers() ([]*model.Member, error)
	GetMember(uuid string) (*model.Member, error)
	GetMemberByUsername(username string) (*model.Member, error)
	GetMemberByEmail(email string) (*model.Member, error)
	CreateMember(body *model.Member) (*model.Member, error)
	UpdateMember(uuid string, body *model.Member) (*model.Member, error)
	UpdateMemberGroups(uuid string, body []*string) error
	DeleteMember(uuid string) (bool, error)
	DeleteMemberByUsername(username string) (bool, error)
	ChangeMemberPassword(uuid string, password string) (bool, error)
}

type MemberAPI struct {
	DB MemberDatabase
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
	member, _ := a.DB.GetMemberByUsername(body.Username)
	if member != nil && member.Username == body.Username && security.CheckPasswordHash(member.Password, body.Password) {
		token, err := security.EncodeJwtToken(member.Username)
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
	q, _ := a.DB.GetMemberByEmail(email)
	message := "email already exists"
	if q == nil {
		message = "email does not exists"
	}
	ResponseHandler(interfaces.Message{Message: message}, nil, ctx)
}

func (a *MemberAPI) CheckUsername(ctx *gin.Context) {
	username := resolveUsername(ctx)
	q, _ := a.DB.GetMemberByUsername(username)
	message := "username already exists"
	if q == nil {
		message = "username does not exists"
	}
	ResponseHandler(interfaces.Message{Message: message}, nil, ctx)
}

func (a *MemberAPI) GetMembers(ctx *gin.Context) {
	q, err := a.DB.GetMembers()
	for _, m := range q {
		m.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) GetMemberByUUID(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.GetMember(uuid)
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
	q, err := a.DB.GetMemberByUsername(username)
	if q != nil {
		q.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) VerifyMember(ctx *gin.Context) {
	username := resolveUsername(ctx)
	member, err := a.DB.GetMemberByUsername(username)
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

func (a *MemberAPI) UpdateMemberGroups(ctx *gin.Context) {
	body, _ := getBodyMemberGroups(ctx)
	uuid := resolveID(ctx)
	err := a.DB.UpdateMemberGroups(uuid, body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(interfaces.Message{Message: "member groups updated successfully"}, err, ctx)
}

func (a *MemberAPI) GetMember(ctx *gin.Context) {
	username := auth.GetAuthorizedUsername(ctx.Request)
	if username == "" {
		ResponseHandler(nil, invalidMemberTokenError, ctx)
		return
	}
	q, err := a.DB.GetMemberByUsername(username)
	if q != nil {
		q.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) UpdateMember(ctx *gin.Context) {
	username := auth.GetAuthorizedUsername(ctx.Request)
	if username == "" {
		ResponseHandler(nil, invalidMemberTokenError, ctx)
		return
	}
	member, err := a.DB.GetMemberByUsername(username)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	body, _ := getBodyMember(ctx)
	body.State = member.State
	q, err := a.DB.UpdateMember(member.UUID, body)
	if q != nil {
		q.MaskPassword()
	}
	ResponseHandler(q, err, ctx)
}

func (a *MemberAPI) DeleteMember(ctx *gin.Context) {
	username := auth.GetAuthorizedUsername(ctx.Request)
	if username == "" {
		ResponseHandler(nil, invalidMemberTokenError, ctx)
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
	username := auth.GetAuthorizedUsername(ctx.Request)
	if username == "" {
		ResponseHandler(nil, invalidMemberTokenError, ctx)
		return
	}
	member, err := a.DB.GetMemberByUsername(username)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	if member.Username != body.Username || !security.CheckPasswordHash(member.Password, body.Password) {
		ResponseHandler(nil, nerrors.NewErrUnauthorized("invalid username or password"), ctx)
		return
	}
	_, err = a.DB.ChangeMemberPassword(member.UUID, body.NewPassword)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(interfaces.Message{Message: "your password has been changed successfully"}, err, ctx)
}

func (a *MemberAPI) RefreshToken(ctx *gin.Context) {
	username := auth.GetAuthorizedUsername(ctx.Request)
	if username == "" {
		ResponseHandler(nil, invalidMemberTokenError, ctx)
		return
	}
	token, err := security.EncodeJwtToken(username)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	ResponseHandler(&model.TokenResponse{AccessToken: token, TokenType: "JWT"}, err, ctx)
}
