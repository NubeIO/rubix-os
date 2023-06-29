package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/args"
	"github.com/gin-gonic/gin"
)

type DeviceDatabase interface {
	GetDevices(args args.Args) ([]*model.Device, error)
	GetDevice(uuid string, args args.Args) (*model.Device, error)
	GetOneDeviceByArgs(args args.Args) (*model.Device, error)
	GetDeviceByName(networkName string, deviceName string, args args.Args) (*model.Device, error)
	CreateDevice(body *model.Device) (*model.Device, error)
	UpdateDevice(uuid string, body *model.Device) (*model.Device, error)
	DeleteDevice(uuid string) (bool, error)
	DeleteOneDeviceByArgs(args args.Args) (bool, error)
	DeleteDeviceByName(networkName string, deviceName string, args args.Args) (bool, error)

	CreateDevicePlugin(body *model.Device) (*model.Device, error)
	UpdateDevicePlugin(uuid string, body *model.Device) (*model.Device, error)
	DeleteDevicePlugin(uuid string) (bool, error)

	CreateDeviceMetaTags(deviceUUID string, deviceMetaTags []*model.DeviceMetaTag) ([]*model.DeviceMetaTag, error)
}
type DeviceAPI struct {
	DB DeviceDatabase
}

func (a *DeviceAPI) GetDevices(ctx *gin.Context) {
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetDevices(args)
	ResponseHandler(q, err, ctx)
}

func (a *DeviceAPI) GetDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetDevice(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *DeviceAPI) GetOneDeviceByArgs(ctx *gin.Context) {
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetOneDeviceByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (a *DeviceAPI) GetDeviceByName(ctx *gin.Context) {
	networkName := resolveNetworkName(ctx)
	deviceName := resolveDeviceName(ctx)
	args := buildDeviceArgs(ctx)
	q, err := a.DB.GetDeviceByName(networkName, deviceName, args)
	ResponseHandler(q, err, ctx)
}

func (a *DeviceAPI) UpdateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateDevicePlugin(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *DeviceAPI) CreateDevice(ctx *gin.Context) {
	body, _ := getBODYDevice(ctx)
	q, err := a.DB.CreateDevicePlugin(body)
	ResponseHandler(q, err, ctx)
}

func (a *DeviceAPI) DeleteDevice(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteDevicePlugin(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *DeviceAPI) DeleteOneDeviceByArgs(ctx *gin.Context) {
	args := buildDeviceArgs(ctx)
	q, err := a.DB.DeleteOneDeviceByArgs(args)
	ResponseHandler(q, err, ctx)
}

func (a *DeviceAPI) DeleteDeviceByName(ctx *gin.Context) {
	networkName := resolveNetworkName(ctx)
	deviceName := resolveDeviceName(ctx)
	args := buildDeviceArgs(ctx)
	q, err := a.DB.DeleteDeviceByName(networkName, deviceName, args)
	ResponseHandler(q, err, ctx)
}

func (a *DeviceAPI) CreateDeviceMetaTags(ctx *gin.Context) {
	deviceUUID := resolveID(ctx)
	body, _ := getBodyBulkDeviceMetaTag(ctx)
	q, err := a.DB.CreateDeviceMetaTags(deviceUUID, body)
	if err != nil {
		ResponseHandler(q, err, ctx)
		return
	}
	ResponseHandler(q, err, ctx)
}
