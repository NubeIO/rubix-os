package api

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/services/appstore"
	"github.com/NubeIO/rubix-os/utils/helpers"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type AppStoreApi struct {
	Store *appstore.Store
}

func (a *AppStoreApi) UploadAddOnAppStore(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	m := &interfaces.Upload{
		Name:    c.Query("name"),
		Version: c.Query("version"),
		Arch:    c.Query("arch"),
		File:    file,
	}
	data, err := a.Store.UploadAddOnAppStore(m)
	ResponseHandler(data, err, c)
}

func (a *AppStoreApi) CheckAppExistence(c *gin.Context) {
	name := c.Query("name")
	arch := c.Query("arch")
	version := c.Query("version")
	if err := a.checkAppExistence(name, arch, version); err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(interfaces.FoundMessage{Found: true}, nil, c)
}

func (a *AppStoreApi) checkAppExistence(name, arch, version string) error {
	if name == "" {
		return errors.New("name can not be empty")
	}
	if err := helpers.CheckVersion(version); err != nil {
		return err
	}
	if arch == "" {
		return errors.New("arch can not be empty")
	}
	p := global.Installer.GetAppsStoreAppPathWithArchVersion(name, arch, version)
	found := fileutils.DirExists(p)
	if !found {
		return errors.New(fmt.Sprintf("failed to find app: %s with arch: %s & version: %s with  in app store", name, arch, version))
	}
	files, _ := ioutil.ReadDir(p)
	if len(files) == 0 {
		return errors.New(fmt.Sprintf("failed to find app: %s with arch: %s & version: %s with  in app store", name, arch, version))
	}
	return nil
}
