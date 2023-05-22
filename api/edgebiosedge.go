package api

import (
	"errors"
	"github.com/NubeIO/flow-framework/global"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/flow-framework/src/cli/cligetter"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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

func (a *EdgeBiosEdgeApi) EdgeBiosRubixEdgeUpload(ctx *gin.Context) {
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
	data, err := cli.RubixEdgeUpload(m)
	ResponseHandler(data, err, ctx)
}

func (a *EdgeBiosEdgeApi) EdgeBiosRubixEdgeInstall(ctx *gin.Context) {
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
	data, err := cli.RubixEdgeInstall(m.Version)
	ResponseHandler(data, err, ctx)
}

func (a *EdgeBiosEdgeApi) EdgeBiosGetRubixEdgeVersion(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	data, err := cli.GetRubixEdgeVersion()
	ResponseHandler(data, err, ctx)
}

func (a *EdgeBiosEdgeApi) attachFileOnModel(m *interfaces.FileUpload) error {
	storePath := global.Installer.GetAppsStoreAppPathWithArchVersion("rubix-edge", m.Arch, m.Version)
	files, err := ioutil.ReadDir(storePath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("rubix-edge store file doesn't exist")
	}
	m.File = path.Join(storePath, files[0].Name())
	return nil
}
