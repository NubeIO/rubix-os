package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
)

type HostDatabase interface {
	GetHosts(withOpenVpn bool, args argspkg.Args) ([]*model.Host, error)
	GetHost(uuid string, args argspkg.Args) (*model.Host, error)
	CreateHost(body *model.Host) (*model.Host, error)
	UpdateHost(uuid string, body *model.Host) (*model.Host, error)
	DeleteHost(uuid string) (*interfaces.Message, error)
	DropHosts() (*interfaces.Message, error)
	ConfigureOpenVPN(uuid string) (*interfaces.Message, error)
}

type HostAPI struct {
	DB HostDatabase
}

func (a *HostAPI) GetHostSchema(ctx *gin.Context) {
	q := interfaces.GetHostSchema()
	ResponseHandler(q, nil, ctx)
}

func (a *HostAPI) GetHosts(ctx *gin.Context) {
	args := buildHostArgs(ctx)
	q, err := a.DB.GetHosts(true, args)
	ResponseHandler(q, err, ctx)
}

func (a *HostAPI) GetHost(ctx *gin.Context) {
	uuid := resolveID(ctx)
	args := buildHostArgs(ctx)
	q, err := a.DB.GetHost(uuid, args)
	ResponseHandler(q, err, ctx)
}

func (a *HostAPI) CreateHost(ctx *gin.Context) {
	body, _ := getBodyHost(ctx)
	q, err := a.DB.CreateHost(body)
	ResponseHandler(q, err, ctx)
}

func (a *HostAPI) UpdateHost(ctx *gin.Context) {
	body, _ := getBodyHost(ctx)
	uuid := resolveID(ctx)
	q, err := a.DB.UpdateHost(uuid, body)
	ResponseHandler(q, err, ctx)
}

func (a *HostAPI) DeleteHost(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.DeleteHost(uuid)
	ResponseHandler(q, err, ctx)
}

func (a *HostAPI) DropHosts(ctx *gin.Context) {
	q, err := a.DB.DropHosts()
	ResponseHandler(q, err, ctx)
}

func (a *HostAPI) ConfigureOpenVPN(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := a.DB.ConfigureOpenVPN(uuid)
	ResponseHandler(q, err, ctx)
}
