package api

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/nerrors"
	"github.com/NubeIO/rubix-os/utils/nstring"
	"github.com/gin-gonic/gin"
)

type MemberDeviceDatabase interface {
	GetMemberDevicesByMemberUUID(memberUUID string) ([]*model.MemberDevice, error)
	GetMemberDevicesByArgs(args args.Args) ([]*model.MemberDevice, error)
	GetOneMemberDeviceByArgs(args args.Args) (*model.MemberDevice, error)
	CreateMemberDevice(body *model.MemberDevice) (*model.MemberDevice, error)
	UpdateMemberDevice(uuid string, body *model.MemberDevice) (*model.MemberDevice, error)
	DeleteMemberDevicesByArgs(args args.Args) (bool, error)

	GetMemberByUsername(username string, args args.Args) (*model.Member, error)
}

type MemberDeviceAPI struct {
	DB MemberDeviceDatabase
}

func (a *MemberDeviceAPI) GetMemberDevices(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	member, err := a.DB.GetMemberByUsername(username, args.Args{})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	q, err := a.DB.GetMemberDevicesByMemberUUID(member.UUID)
	ResponseHandler(q, err, ctx)
}

func (a *MemberDeviceAPI) GetMemberDevice(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	member, err := a.DB.GetMemberByUsername(username, args.Args{})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	deviceId := resolveDeviceId(ctx)
	q, err := a.DB.GetMemberDevicesByArgs(args.Args{DeviceId: nstring.New(deviceId), MemberUUID: nstring.New(member.UUID)})
	ResponseHandler(q, err, ctx)
}

func (a *MemberDeviceAPI) CreateMemberDevice(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	member, err := a.DB.GetMemberByUsername(username, args.Args{})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	body, _ := getBodyMemberDevice(ctx)
	args := args.Args{DeviceId: nstring.New(body.DeviceID), MemberUUID: nstring.New(member.UUID)}
	memberDevice, err := a.DB.GetOneMemberDeviceByArgs(args)
	if memberDevice != nil {
		ResponseHandler(memberDevice, err, ctx)
		return
	}
	body.MemberUUID = member.UUID
	q, err := a.DB.CreateMemberDevice(body)
	ResponseHandler(q, err, ctx)
}

func (a *MemberDeviceAPI) UpdateMemberDevice(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	member, err := a.DB.GetMemberByUsername(username, args.Args{})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	body, _ := getBodyMemberDevice(ctx)
	deviceId := resolveDeviceId(ctx)
	args := args.Args{DeviceId: &deviceId, MemberUUID: nstring.New(member.UUID)}
	memberDevice, err := a.DB.GetOneMemberDeviceByArgs(args)
	if memberDevice == nil {
		ResponseHandler(nil, errors.New("device not found"), ctx)
		return
	}
	q, err := a.DB.UpdateMemberDevice(memberDevice.UUID, body)
	ResponseHandler(q, err, ctx)
}

func (a *MemberDeviceAPI) DeleteMemberDevice(ctx *gin.Context) {
	username, err := getAuthorizedUsername(ctx.Request)
	if err != nil {
		ResponseHandler(nil, nerrors.NewErrUnauthorized(err.Error()), ctx)
		return
	}
	member, err := a.DB.GetMemberByUsername(username, args.Args{})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	deviceId := resolveDeviceId(ctx)
	q, err := a.DB.DeleteMemberDevicesByArgs(args.Args{DeviceId: nstring.New(deviceId), MemberUUID: nstring.New(member.UUID)})
	ResponseHandler(q, err, ctx)
}
