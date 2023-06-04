package api

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/src/cli/cligetter"
	"github.com/NubeIO/rubix-os/src/cli/constants"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"path"
)

type EdgeBiosEdgeDatabase interface {
	ResolveHost(uuid string, name string) (*model.Host, error)
}

type EdgeBiosEdgeApi struct {
	DB EdgeBiosEdgeDatabase
}

func (a *EdgeBiosEdgeApi) EdgeBiosRubixOsUpload(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	var m *interfaces.FileUpload
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	err = a.attachFileOnModel(m)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	data, err := cli.RubixOsUpload(m)
	ResponseHandler(data, err, ctx)
}

func (a *EdgeBiosEdgeApi) EdgeBiosRubixOsInstall(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	var m *interfaces.FileUpload
	err = ctx.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	data, err := cli.RubixOsInstall(m.Version)
	ResponseHandler(data, err, ctx)
}

func (a *EdgeBiosEdgeApi) EdgeBiosGetRubixOsVersion(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	data, err := cli.GetEdgeRubixOsVersion()
	ResponseHandler(data, err, ctx)
}

func (a *EdgeBiosEdgeApi) attachFileOnModel(m *interfaces.FileUpload) error {
	storePath := global.Installer.GetAppsStoreAppPathWithArchVersion(constants.RubixOs, m.Arch, m.Version)
	files, err := ioutil.ReadDir(storePath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("app store file doesn't exist")
	}
	m.File = path.Join(storePath, files[0].Name())
	return nil
}
